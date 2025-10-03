package files

import (
    "os"
    "path/filepath"

    "github.com/e74000/xtd/set"
)

func Walk(root string) (files *set.Set[string], dirs *set.Set[string], err error) {
    abs, err := filepath.Abs(root)
    if err != nil {
        return nil, nil, err
    }

    files = set.New[string]()
    dirs = set.New[string]()

    err = filepath.WalkDir(abs, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }

        rel, err := filepath.Rel(abs, path)
        if err != nil {
            return err
        }

        if d.IsDir() {
            dirs.Add(rel)
        } else {
            files.Add(rel)
        }

        return nil
    })

    return files, dirs, err
}

func WalkFiles(root string) (files *set.Set[string], err error) {
    abs, err := filepath.Abs(root)
    if err != nil {
        return nil, err
    }

    files = set.New[string]()

    err = filepath.WalkDir(abs, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if !d.IsDir() {
            files.Add(path)
        }

        return nil
    })

    return files, err
}

func WalkDirs(root string) (dirs *set.Set[string], err error) {
    abs, err := filepath.Abs(root)
    if err != nil {
        return nil, err
    }

    dirs = set.New[string]()

    err = filepath.WalkDir(abs, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if d.IsDir() {
            dirs.Add(path)
        }

        return nil
    })

    return dirs, err
}

func WalkInfo(root string) (files map[string]os.FileInfo, dirs map[string]os.FileInfo, err error) {
    abs, err := filepath.Abs(root)
    if err != nil {
        return nil, nil, err
    }

    files = make(map[string]os.FileInfo)
    dirs = make(map[string]os.FileInfo)

    err = filepath.WalkDir(abs, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }

        rel, err := filepath.Rel(abs, path)
        if err != nil {
            return err
        }

        info, err := d.Info()
        if err != nil {
            return err
        }

        if d.IsDir() {
            dirs[rel] = info
        } else {
            files[rel] = info
        }

        return nil
    })

    return files, dirs, err
}

func WalkFilesInfo(root string) (files map[string]os.FileInfo, err error) {
    abs, err := filepath.Abs(root)
    if err != nil {
        return nil, err
    }

    files = make(map[string]os.FileInfo)

    err = filepath.WalkDir(abs, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }

        rel, err := filepath.Rel(abs, path)
        if err != nil {
            return err
        }

        info, err := d.Info()
        if err != nil {
            return err
        }

        if !d.IsDir() {
            files[rel] = info
        }

        return nil
    })

    return files, err
}

func WalkDirsInfo(root string) (dirs map[string]os.FileInfo, err error) {
    abs, err := filepath.Abs(root)
    if err != nil {
        return nil, err
    }

    dirs = make(map[string]os.FileInfo)

    err = filepath.WalkDir(abs, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }

        rel, err := filepath.Rel(abs, path)
        if err != nil {
            return err
        }

        info, err := d.Info()
        if err != nil {
            return err
        }

        if d.IsDir() {
            dirs[rel] = info
        }

        return nil
    })

    return dirs, err
}
