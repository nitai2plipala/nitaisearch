package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "runtime"
    "strings"
)

type Bookmark struct {
    Name     string     `json:"name"`
    URL      string     `json:"url"`
    Children []Bookmark `json:"children"`
}

type ChromeBookmarksRoot struct {
    Roots struct {
        BookmarkBar Bookmark `json:"bookmark_bar"`
        Other       Bookmark `json:"other"`
        Synced      Bookmark `json:"synced"`
    } `json:"roots"`
}

func printBookmarks(b Bookmark, prefix string) {
    if b.URL != "" {
        fmt.Printf("%s%s -> %s\n", prefix, b.Name, b.URL)
    }
    for _, child := range b.Children {
        printBookmarks(child, prefix+"  ")
    }
}

func getBrowserProfiles(browser string) ([]string, error) {
    home, _ := os.UserHomeDir()
    var base string
    switch runtime.GOOS {
    case "windows":
        switch browser {
        case "chrome":
            base = filepath.Join(home, "AppData", "Local", "Google", "Chrome", "User Data")
        case "edge":
            base = filepath.Join(home, "AppData", "Local", "Microsoft", "Edge", "User Data")
        case "brave":
            base = filepath.Join(home, "AppData", "Local", "BraveSoftware", "Brave-Browser", "User Data")
        }
    case "darwin":
        switch browser {
        case "chrome":
            base = filepath.Join(home, "Library", "Application Support", "Google", "Chrome")
        case "edge":
            base = filepath.Join(home, "Library", "Application Support", "Microsoft Edge")
        case "brave":
            base = filepath.Join(home, "Library", "Application Support", "BraveSoftware", "Brave-Browser")
        }
    case "linux":
        switch browser {
        case "chrome":
            base = filepath.Join(home, ".config", "google-chrome")
        case "edge":
            base = filepath.Join(home, ".config", "microsoft-edge")
        case "brave":
            base = filepath.Join(home, ".config", "BraveSoftware", "Brave-Browser")
        }
    }

    dirs, err := ioutil.ReadDir(base)
    if err != nil {
        return nil, err
    }

    var profiles []string
    for _, d := range dirs {
        if d.IsDir() && (d.Name() == "Default" || strings.HasPrefix(d.Name(), "Profile")) {
            profiles = append(profiles, d.Name())
        }
    }
    return profiles, nil
}

func main() {
    browsers := []string{"chrome", "edge", "brave"}
    fmt.Println("Select browser:")
    for i, b := range browsers {
        fmt.Printf("%d) %s\n", i+1, b)
    }

    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter choice (number): ")
    choiceStr, _ := reader.ReadString('\n')
    choiceStr = strings.TrimSpace(choiceStr)
    choice := 0
    fmt.Sscanf(choiceStr, "%d", &choice)
    if choice < 1 || choice > len(browsers) {
        fmt.Println("Invalid choice")
        return
    }

    browser := browsers[choice-1]
    profiles, err := getBrowserProfiles(browser)
    if err != nil || len(profiles) == 0 {
        fmt.Println("Cannot find any profiles")
        return
    }

    fmt.Println("Select profile:")
    for i, p := range profiles {
        fmt.Printf("%d) %s\n", i+1, p)
    }
    fmt.Print("Enter choice (number): ")
    profileChoiceStr, _ := reader.ReadString('\n')
    profileChoiceStr = strings.TrimSpace(profileChoiceStr)
    profileChoice := 0
    fmt.Sscanf(profileChoiceStr, "%d", &profileChoice)
    if profileChoice < 1 || profileChoice > len(profiles) {
        fmt.Println("Invalid choice")
        return
    }

    profile := profiles[profileChoice-1]
    home, _ := os.UserHomeDir()
    var bookmarkPath string
    switch runtime.GOOS {
    case "windows":
        bookmarkPath = filepath.Join(home, "AppData", "Local", "Google", "Chrome", "User Data", profile, "Bookmarks")
    case "darwin":
        bookmarkPath = filepath.Join(home, "Library", "Application Support", "Google", "Chrome", profile, "Bookmarks")
    case "linux":
        bookmarkPath = filepath.Join(home, ".config", "google-chrome", profile, "Bookmarks")
    }

    data, err := ioutil.ReadFile(bookmarkPath)
    if err != nil {
        fmt.Println("Error reading bookmarks:", err)
        return
    }

    var root ChromeBookmarksRoot
    if err := json.Unmarshal(data, &root); err != nil {
        fmt.Println("Error parsing JSON:", err)
        return
    }

    fmt.Println("=== Bookmark Bar ===")
    printBookmarks(root.Roots.BookmarkBar, "")
    fmt.Println("\n=== Other Bookmarks ===")
    printBookmarks(root.Roots.Other, "")
    fmt.Println("\n=== Synced Bookmarks ===")
    printBookmarks(root.Roots.Synced, "")
}
