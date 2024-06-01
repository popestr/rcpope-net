/**
 * Represents a single course.
 * @typedef {Object} Course
 * @property {string} course_code - The code for the course.
 * @property {string} course_topic - The course topic.
 * @property {string} course_name - The name of the course.
 * @property {string} course_classes - A space-separated list of classifications for styling.
 * @property {string} semester - The semester in which the course is taught.
 * @property {string} classification - The course classification type.
 * @property {string} languages - The programming languages related to the course.
 * @property {string} classification_icons - The HTML or text representing icons for classifications.
 * @property {string} lang_icons - The HTML or text representing icons for languages.
 */

/**
 * Represents a single classification.
 * @typedef {Object} Classification
 * @property {string} abbreviation - The short form of the classification.
 * @property {string} icon_html - The HTML representation of the classification icon.
 * @property {string} longname - The full name of the classification.
 */

/**
 * Represents a single language.
 * @typedef {Object} Language
 * @property {string} abbreviation - The short form of the language.
 * @property {string} icon_html - The HTML representation of the language icon.
 * @property {string} longname - The full name of the language.
 */

/**
 * The main structure of the courses data.
 * @typedef {Object} CoursesData
 * @property {Array.<Course>} courses - An array of courses.
 * @property {Array.<Classification>} classifications - An array of classifications.
 * @property {Array.<Language>} languages - An array of languages.
 */
