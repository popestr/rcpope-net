/**
 * Represents a single course.
 * @typedef {Object} Course
 * @property {string} course_code - The code for the course.
 * @property {string} course_topic - The course topic.
 * @property {string} course_name - The name of the course.
 * @property {string} course_classes - A space-separated list of classifications for styling.
 * @property {string} semester - The semester in which the course is taught.
 * @property {string[]} classifications - The course classification types.
 * @property {string[]} languages - The programming languages related to the course.
 */

/**
 * Represents a single classification.
 * @typedef {Object} Classification
 * @property {string} abbreviation - The short form of the classification.
 * @property {string} longname - The full name of the classification.
 */

/**
 * Represents a single language.
 * @typedef {Object} Language
 * @property {string} abbreviation - The short form of the language.
 * @property {string} longname - The full name of the language.
 */

/**
 * The main structure of the courses data.
 * @typedef {Object} CoursesData
 * @property {Array.<Course>} courses - An array of courses.
 * @property {Array.<Classification>} classifications - An array of classifications.
 * @property {Array.<Language>} languages - An array of languages.
 */
