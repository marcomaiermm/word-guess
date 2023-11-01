function setupTheme() {
  // On page load or when changing themes, best to add inline in `head` to avoid FOUC
  if (
    localStorage.theme === "dark" ||
    (!("theme" in localStorage) &&
      window.matchMedia("(prefers-color-scheme: dark)").matches)
  ) {
    document.documentElement.classList.add("dark");
  } else {
    document.documentElement.classList.remove("dark");
  }
}

function switchTheme(theme: string) {
  if (theme === "dark" || theme === "light") {
    localStorage.theme = theme;
  } else {
    localStorage.removeItem("theme");
  }
  setupTheme();
}

setupTheme();
