package files

import (
    "errors"
    "io"
    "os"
    "path/filepath"
    "slices"
)

func AtomicWrite(path string, gen func(w io.Writer) error) error {
    dir, base := filepath.Split(path)

    tmp, err := os.CreateTemp(dir, "."+base+".tmp-*")
    if err != nil {
        return err
    }
    defer func(tmp *os.File) {
        _ = tmp.Close()
        _ = os.Remove(tmp.Name())
    }(tmp)

    if err := gen(tmp); err != nil {
        return err
    }
    if err := tmp.Sync(); err != nil {
        return err
    }
    if err := tmp.Close(); err != nil {
        return err
    }

    if err := os.Rename(tmp.Name(), path); err != nil {
        return err
    }

    if df, err := os.Open(dir); err == nil {
        _ = df.Sync()
        _ = df.Close()
    }
    return nil
}

func AtomicEdit(path string, gen func(w io.Writer) error) error {
    dir, base := filepath.Split(path)

    tmp, err := os.CreateTemp(dir, "."+base+".tmp-*")
    if err != nil {
        return err
    }
    defer func(tmp *os.File) {
        _ = tmp.Close()
        _ = os.Remove(tmp.Name())
    }(tmp)

    if err := gen(tmp); err != nil {
        return err
    }
    if err := tmp.Sync(); err != nil {
        return err
    }
    if err := tmp.Close(); err != nil {
        return err
    }

    if eq, err := cmpFiles(tmp.Name(), path); err != nil {
        return err
    } else if eq {
        return nil
    }

    if err := os.Rename(tmp.Name(), path); err != nil {
        return err
    }

    if df, err := os.Open(dir); err == nil {
        _ = df.Sync()
        _ = df.Close()
    }

    return nil
}

func cmpFiles(a, b string) (bool, error) {
    aFi, err := os.Stat(a)
    if err != nil {
        return false, err
    } else if aFi.IsDir() {
        return false, errors.New("a is a directory")
    }

    bFi, err := os.Stat(b)
    if err != nil {
        return false, err
    } else if bFi.IsDir() {
        return false, errors.New("b is a directory")
    }

    if aFi.Size() != bFi.Size() {
        return false, nil
    }

    return filesEqualChunked(a, b)
}

func filesEqualChunked(a, b string) (bool, error) {
    aF, err := os.Open(a)
    if err != nil {
        return false, err
    }
    // noinspection GoUnhandledErrorResult
    defer aF.Close()

    bF, err := os.Open(b)
    if err != nil {
        return false, err
    }
    // noinspection GoUnhandledErrorResult
    defer bF.Close()

    const bufSize = 128 * 1024
    aBuf := make([]byte, bufSize)
    bBuf := make([]byte, bufSize)

    for {
        aN, aErr := io.ReadFull(aF, aBuf)
        bN, bErr := io.ReadFull(bF, bBuf)

        if aErr != nil && !errors.Is(aErr, io.EOF) && !errors.Is(aErr, io.ErrUnexpectedEOF) {
            return false, aErr
        }

        if bErr != nil && !errors.Is(bErr, io.EOF) && !errors.Is(bErr, io.ErrUnexpectedEOF) {
            return false, bErr
        }

        if aN == 0 && bN == 0 {
            return true, nil
        }

        if aN != bN || !slices.Equal(aBuf[:aN], bBuf[:bN]) {
            return false, nil
        }
    }
}
