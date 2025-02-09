document.addEventListener("DOMContentLoaded", function () {
  const lightIcon = "fa-sun";
  const darkIcon = "fa-moon";

  let isDark =
    window.matchMedia &&
    window.matchMedia("(prefers-color-scheme: dark)").matches;

  const siteNav = document.querySelector("nav.site-nav");
  const icon = document.createElement("i");
  icon.classList.add("fas");
  const span = document.createElement("span");
  span.innerHTML = "Toggle dark mode";
  const link = document.createElement("a");
  link.classList.add("nav-list-link");
  link.append(icon, span);
  link.setAttribute("href", "#");
  const li = document.createElement("li");
  li.classList.add("nav-list-item");
  li.append(link);
  const ul = document.createElement("ul");
  ul.classList.add("nav-list", "nav-list-site-settings");
  ul.append(li);
  siteNav.append(ul);

  function update() {
    if (isDark) {
      jtd.setTheme("dark");
      icon.classList.remove(darkIcon);
      icon.classList.add(lightIcon);
    } else {
      jtd.setTheme("light");
      icon.classList.remove(lightIcon);
      icon.classList.add(darkIcon);
    }
  }

  update();

  function toggle() {
    isDark = !isDark;
    update();
  }

  jtd.addEvent(link, "click", function (event) {
    event.preventDefault();
    toggle();
  });
});
