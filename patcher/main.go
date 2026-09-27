package main

import (
    "bufio"
    "bytes"
    "compress/zlib"
    "crypto/sha256"
    _ "embed"
    "encoding/base64"
    "encoding/binary"
    "encoding/hex"
    "errors"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
)

const (
    BuildName = "R1"
    GameExeName = "hedge.exe"
    RetailSHA256 = "81ce80f1bd5cc74183f621871e3ec4fa079bf0d694b652f7b20002cea82c1f21"
    TargetSHA256 = "9f6cf822e1927c4968dc22cc4328ea86cdd658cf486843498208be782c16dcfe"
    RetailSize int64 = 3124397
    TargetSize int64 = 3124397
)

//go:embed over_the_hedge_r1.rtdp1.zlib.b64
var embeddedPatchB64 string

type options struct { gameDir string; noPause bool }

func parseArgs() options {
    var o options
    args := os.Args[1:]
    for i := 0; i < len(args); i++ {
        switch args[i] {
        case "--game-dir":
            if i+1 < len(args) { i++; o.gameDir = args[i] }
        case "--no-pause":
            o.noPause = true
        }
    }
    return o
}

func pause(noPause bool) {
    if noPause { return }
    fmt.Print("\nPress Enter to close...")
    _, _ = bufio.NewReader(os.Stdin).ReadString('\n')
}

func sha256Bytes(b []byte) string {
    s := sha256.Sum256(b)
    return hex.EncodeToString(s[:])
}

func sha256File(p string) (string, error) {
    f, e := os.Open(p)
    if e != nil { return "", e }
    defer f.Close()
    h := sha256.New()
    if _, e = io.Copy(h, f); e != nil { return "", e }
    return hex.EncodeToString(h.Sum(nil)), nil
}

func appDir(o options) (string, error) {
    if o.gameDir != "" { return filepath.Abs(o.gameDir) }
    e, err := os.Executable()
    if err != nil { return "", err }
    return filepath.Dir(e), nil
}

func copyFile(src, dst string) error {
    in, e := os.Open(src)
    if e != nil { return e }
    defer in.Close()
    st, e := in.Stat()
    if e != nil { return e }

    out, e := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, st.Mode())
    if e != nil { return e }
    ok := false
    defer func() {
        _ = out.Close()
        if !ok { _ = os.Remove(dst) }
    }()

    if _, e = io.Copy(out, in); e != nil { return e }
    if e = out.Sync(); e != nil { return e }
    if e = out.Close(); e != nil { return e }
    ok = true
    return nil
}

func backupRetail(exe string) (string, error) {
    for _, c := range []string{exe + ".Backup", exe + ".Backup.Vanilla"} {
        if _, e := os.Stat(c); e == nil {
            h, e := sha256File(c)
            if e == nil && strings.EqualFold(h, RetailSHA256) { return c, nil }
            continue
        } else if !os.IsNotExist(e) {
            return "", e
        }

        if e := copyFile(exe, c); e != nil { return "", e }
        h, e := sha256File(c)
        if e != nil { return "", e }
        if !strings.EqualFold(h, RetailSHA256) {
            _ = os.Remove(c)
            return "", errors.New("backup verification failed")
        }
        return c, nil
    }
    return "", errors.New("existing backup names are occupied by non-original files")
}

func applyRTDP1(source, compressed []byte) ([]byte, error) {
    zr, e := zlib.NewReader(bytes.NewReader(compressed))
    if e != nil { return nil, fmt.Errorf("invalid compressed patch: %w", e) }
    raw, e := io.ReadAll(zr)
    if ce := zr.Close(); e == nil { e = ce }
    if e != nil { return nil, e }

    r := bytes.NewReader(raw)
    magic := make([]byte, 6)
    if _, e = io.ReadFull(r, magic); e != nil || string(magic) != "RTDP1\x00" {
        return nil, errors.New("invalid patch format")
    }

    var ss, ts uint64
    if e = binary.Read(r, binary.LittleEndian, &ss); e != nil { return nil, errors.New("truncated patch header") }
    if e = binary.Read(r, binary.LittleEndian, &ts); e != nil { return nil, errors.New("truncated patch header") }

    sh := make([]byte, 32)
    th := make([]byte, 32)
    if _, e = io.ReadFull(r, sh); e != nil { return nil, errors.New("truncated source hash") }
    if _, e = io.ReadFull(r, th); e != nil { return nil, errors.New("truncated target hash") }

    var count uint32
    if e = binary.Read(r, binary.LittleEndian, &count); e != nil { return nil, errors.New("truncated instruction count") }

    a := sha256.Sum256(source)
    if uint64(len(source)) != ss || !bytes.Equal(a[:], sh) {
        return nil, errors.New("patch source verification failed")
    }

    out := bytes.NewBuffer(make([]byte, 0, int(ts)))
    for i := uint32(0); i < count; i++ {
        op, e := r.ReadByte()
        if e != nil { return nil, errors.New("truncated patch instruction stream") }
        switch op {
        case 0:
            var off uint64
            var ln uint32
            if e = binary.Read(r, binary.LittleEndian, &off); e != nil { return nil, e }
            if e = binary.Read(r, binary.LittleEndian, &ln); e != nil { return nil, e }
            end := off + uint64(ln)
            if end > uint64(len(source)) { return nil, errors.New("copy range outside source") }
            out.Write(source[off:end])
        case 1:
            var ln uint32
            if e = binary.Read(r, binary.LittleEndian, &ln); e != nil { return nil, e }
            lit := make([]byte, int(ln))
            if _, e = io.ReadFull(r, lit); e != nil { return nil, errors.New("truncated literal data") }
            out.Write(lit)
        default:
            return nil, fmt.Errorf("unknown patch opcode %d", op)
        }
    }

    if r.Len() != 0 { return nil, errors.New("unexpected trailing patch data") }
    result := out.Bytes()
    t := sha256.Sum256(result)
    if uint64(len(result)) != ts || !bytes.Equal(t[:], th) {
        return nil, errors.New("patch target verification failed")
    }
    return result, nil
}

func installVerified(exe string, target []byte) error {
    dir := filepath.Dir(exe)
    stage := filepath.Join(dir, GameExeName+".PatchNew")
    rollback := filepath.Join(dir, GameExeName+".RollbackTemp")
    _ = os.Remove(stage)
    _ = os.Remove(rollback)

    if e := os.WriteFile(stage, target, 0755); e != nil { return e }
    h, e := sha256File(stage)
    if e != nil || !strings.EqualFold(h, TargetSHA256) {
        _ = os.Remove(stage)
        return errors.New("staged R1 verification failed")
    }

    if e = os.Rename(exe, rollback); e != nil {
        _ = os.Remove(stage)
        return fmt.Errorf("cannot prepare rollback file: %w", e)
    }
    if e = os.Rename(stage, exe); e != nil {
        _ = os.Rename(rollback, exe)
        _ = os.Remove(stage)
        return fmt.Errorf("cannot install patched executable: %w", e)
    }

    h, e = sha256File(exe)
    if e != nil || !strings.EqualFold(h, TargetSHA256) {
        _ = os.Remove(exe)
        _ = os.Rename(rollback, exe)
        return errors.New("final R1 hash mismatch; original restored")
    }

    _ = os.Remove(rollback)
    return nil
}

func main() {
    o := parseArgs()
    exit := 0
    defer func() {
        pause(o.noPause)
        os.Exit(exit)
    }()

    fmt.Println(strings.Repeat("=", 68))
    fmt.Println(" Over the Hedge - Enhanced PC Patch")
    fmt.Println(" Release: R1 | Windows 11 modernization")
    fmt.Println(strings.Repeat("=", 68))

    root, e := appDir(o)
    if e != nil {
        fmt.Println("[ERROR] Cannot determine patcher directory:", e)
        exit = 1
        return
    }
    exe := filepath.Join(root, GameExeName)
    fmt.Println("[OK] Game directory:", root)

    if st, e := os.Stat(exe); e != nil || st.IsDir() {
        fmt.Println("[ERROR] hedge.exe not found next to the patcher.")
        fmt.Println("        No files have been modified.")
        exit = 2
        return
    }

    fmt.Println("[ .. ] Verifying game executable...")
    h, e := sha256File(exe)
    if e != nil {
        fmt.Println("[ERROR] Cannot hash hedge.exe:", e)
        exit = 1
        return
    }

    if strings.EqualFold(h, TargetSHA256) {
        fmt.Println("[OK] R1 is already installed. Nothing to do.")
        return
    }

    if !strings.EqualFold(h, RetailSHA256) {
        fmt.Println("[ERROR] Unsupported hedge.exe.")
        fmt.Println("        The patcher supports only the unpacked/deprotected original PC executable.")
        fmt.Println("        No files have been modified.")
        fmt.Println("        SHA-256:", h)
        exit = 2
        return
    }

    st, _ := os.Stat(exe)
    if st.Size() != RetailSize {
        fmt.Println("[ERROR] Original size check failed. No files have been modified.")
        exit = 2
        return
    }

    fmt.Println("[OK] Supported original unpacked executable detected.")
    backup, e := backupRetail(exe)
    if e != nil {
        fmt.Println("[ERROR] Could not create verified backup:", e)
        exit = 1
        return
    }
    fmt.Println("[OK] Backup:", filepath.Base(backup))

    source, e := os.ReadFile(exe)
    if e != nil {
        fmt.Println("[ERROR] Cannot read source executable:", e)
        exit = 1
        return
    }

    fmt.Println("[ .. ] Building cumulative R1...")
    patchData, e := base64.StdEncoding.DecodeString(strings.TrimSpace(embeddedPatchB64))
    if e != nil {
        fmt.Println("[ERROR] Embedded R1 delta is invalid:", e)
        exit = 1
        return
    }
    target, e := applyRTDP1(source, patchData)
    if e != nil {
        fmt.Println("[ERROR]", e)
        fmt.Println("No unverified output was installed.")
        exit = 1
        return
    }

    if int64(len(target)) != TargetSize || !strings.EqualFold(sha256Bytes(target), TargetSHA256) {
        fmt.Println("[ERROR] Rebuilt R1 failed in-memory verification.")
        exit = 1
        return
    }
    fmt.Println("[OK] R1 target hash verified.")

    if e = installVerified(exe, target); e != nil {
        fmt.Println("[ERROR]", e)
        fmt.Println("No unverified patched output was intentionally left installed.")
        exit = 1
        return
    }

    fmt.Println(strings.Repeat("=", 68))
    fmt.Println("[SUCCESS] Over the Hedge Enhanced PC Patch R1 installed.")
    fmt.Println("          Final SHA-256 verification passed.")
    fmt.Println(strings.Repeat("=", 68))
}
