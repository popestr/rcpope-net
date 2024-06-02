const langIcons = {
  angular: "fa-brands fa-angular",
  aws: "fa-brands fa-aws",
  bootstrap: "fa-brands fa-bootstrap",
  cpp: "fa-solid fa-c",
  css: "fa-brands fa-css3-alt",
  golang: "fa-brands fa-golang",
  html: "fa-brands fa-html5",
  java: "fab fa-java",
  js: "fa-brands fa-js-square",
  php: "fa-brands fa-php",
  python: "fab fa-python",
  sql: "fa-solid fa-database",
};

const classIcons = {
  bus: "fas fa-briefcase",
  eng: "fas fa-wrench",
  imp: "fas fa-laptop-code",
  math: "fas fa-square-root-alt",
  musi: "fas fa-music",
  sci: "fas fa-flask",
  stat: "fas fa-chart-bar",
  sts: "fas fa-gavel",
  theo: "fas fa-brain",
};

/**s
 * @param {string} icon
 * @returns {string} The HTML for the icon.
 */
const getIconHtml = (icon) => {
  return `<i class="${icon}"></i>`;
};

/**
 * @param {string} abbreviation The abbreviation of the language or classification.
 * @returns {string} The icon for the abbreviation.
 */
const getIconByAbbreviation = (abbreviation) => {
  if (abbreviation in langIcons) {
    return getIconHtml(langIcons[abbreviation]);
  } else if (abbreviation in classIcons) {
    return getIconHtml(classIcons[abbreviation]);
  }
  return "";
};

/**
 * @param {string[]} abbreviations An array of abbreviations.
 * @returns {string} A string of icons.
 */
const getIconsByAbbreviations = (abbreviations) => {
  if (!abbreviations) {
    return "";
  }
  return abbreviations.map((abbreviation) => getIconByAbbreviation(abbreviation)).join("");
};
