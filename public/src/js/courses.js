/*
    courses.js
    Authored by Ryan Pope <ryan@rcpope.net>
    Licensed under the GNU Public License, v3.0
*/

const SEMESTER_PLACEHOLDER = 'Fall 1969';

/**
 * Sets up the event listeners for course tiles on the Courses page.
 * 
 * @param {string} courseSelector The selector matching all course tiles.
 */
const initCourseEventListeners = (courseSelector) => {
    document.querySelectorAll(courseSelector).forEach((element) => {
        element.addEventListener("mouseenter", () => {
            const classes = element.getAttribute("data-classifications") || ""; // data-classifications is set by php ajax request handler
            for (const classification of classes.split(" ")) {
                document.querySelector(`.filter-item--${classification}`)?.classList.add("filter-active");
            }
        });

        element.addEventListener("mouseleave", () => {
            document.querySelectorAll(".filter-item").forEach((element) => {
                element.classList.remove("filter-active");
            })
        });
    });
}

/**
 * Toggles a filter in the Courses menu.
 * 
 * @param {Element} element The filter element to toggle.
 * @param {string} classification The classification toggled by the filter.
 */
const toggleFilter = (element, classification) => {
    // check if we are selecting the filter or deselecting it by checking the existing classes
    const delta = element.classList.contains("filter-selected") ? -1 : 1;
    element.classList.toggle("filter-selected");

    // increment result count for filter container (total number of filters selected)
    element.parentElement.setAttribute("data-selected", parseInt(element.parentElement.getAttribute("data-selected")) + delta);

    document.querySelectorAll(".course--" + classification).forEach((node) => {
        // increment match count (stored in data-selectedby) for each course
        const selectedBy = parseInt(node.getAttribute("data-selectedby")) + delta;
        node.setAttribute("data-selectedby", selectedBy);

        // increment result count (stored in data-contains) for parent (.course-group)
        const parentCount = parseInt(node.parentElement?.getAttribute("data-contains")) + delta;
        node.parentElement?.setAttribute("data-contains", parentCount);

        // increment result count (stored in data-contains) for results container (total number of results)
        const mainList = document.querySelector("#main-course-list");
        const totalCount = parseInt(mainList?.getAttribute("data-contains")) + delta;
        mainList?.setAttribute("data-contains", totalCount);
    });
}

/**
 * Sets up the event listeners for filter pills on the Courses page.
 * 
 * @param {string} filterSelector The selector matching all filter pills.
 */
const initFilterEventListeners = (filterSelector) => {
    document.querySelectorAll(filterSelector).forEach((element) => {
        const classification = element.getAttribute("data-classification"); // data-classification is set by php ajax request handler

        // filter hover highlights courses
        element.addEventListener("mouseenter", () => {
            document.querySelectorAll(`.course--${classification}`).forEach((node) => {
                node.classList.add("course-active");
            });
        });

        element.addEventListener("mouseleave", () => {
            document.querySelectorAll(".course").forEach((node) => {
                node.classList.remove("course-active");
            });
        });

        // filter select
        element.onclick = () => toggleFilter(element, classification);
    });
}

/**
 * Collapses/uncollapses a course group header.
 * 
 * @param {Element} groupHeaderElement 
 */
const toggleCourseGroup = (groupHeaderElement) => {
    groupHeaderElement.parentElement?.classList.toggle("course-group--collapsed");
}

/**
 * Initializes course data on the page.
 */
window.onload = () => {
    fetchCourses().then(data => {
        populateCourses(data);
        populateFilters(data);

        initCourseEventListeners(".course");
        initFilterEventListeners(".filter-item");

        setDefaultFilters();
    });
};

/**
 * Fetches courses from the provided API.
 * 
 * @returns {Promise<CoursesData>} A promise that resolves to the courses data.
 */
const fetchCourses = () => {
    return fetch("https://api.rcpope.net/courses_json.php")
        .then(response => {
            if (!response.ok) {
                throw new Error(`Failed to fetch courses. (got ${response.status})`);
            }
            return response.json();
        });
}

/**
 * Populates the courses on the page.
 * 
 * @param {CoursesData} data - The courses data.
 */
const populateCourses = (data) => {
    let prevSemester = SEMESTER_PLACEHOLDER;
    let group = -1;

    for (const course of data.courses) {
        if (course.semester !== prevSemester) {
            group++;
            appendCourseGroupHeader(course.semester, group);
            prevSemester = course.semester;
        }
        appendCourseToGroup(course, group);
    }
}

/**
 * Appends a course group header to the courses list.
 * 
 * @param {string} semester - The semester string.
 * @param {number} group - The group identifier.
 */
const appendCourseGroupHeader = (semester, group) => {
    const groupHeader = document.createElement('div');
    groupHeader.classList.add('course-group-header');
    groupHeader.innerHTML = `<i class="fas fa-caret-down course-group-arrow"></i> ${semester}`;
    groupHeader.onclick = () => toggleCourseGroup(groupHeader);

    const courseGroup = document.createElement('div');
    courseGroup.classList.add('course-group');
    courseGroup.setAttribute('data-contains', '0');
    courseGroup.setAttribute('data-group', group);
    courseGroup.appendChild(groupHeader);

    document.querySelector('.courses-list').appendChild(courseGroup);
}

/**
 * Appends a course to a specific group.
 * 
 * @param {Course} course - The course object.
 * @param {number} group - The group identifier.
 */
const appendCourseToGroup = (course, group) => {
    const courseElement = document.createElement('div');
    courseElement.classList.add('course', ...course.course_classes.split(' '));
    courseElement.setAttribute('data-classifications', `${course.classification} ${course.languages}`);
    courseElement.setAttribute('data-selectedby', '0');

    const courseTitleElement = document.createElement('span');
    courseTitleElement.classList.add('course-title');

    const courseCodeElement = document.createElement('span');
    courseCodeElement.classList.add('course-code');
    courseCodeElement.textContent = course.course_code;

    const courseNameElement = document.createElement('span');
    courseNameElement.classList.add('course-name');
    courseNameElement.textContent = course.course_topic || course.course_name;

    courseTitleElement.appendChild(courseCodeElement);
    courseTitleElement.appendChild(courseNameElement);
    courseElement.appendChild(courseTitleElement);

    const courseIconsElement = document.createElement('span');
    courseIconsElement.classList.add('course-icons');

    const courseClassificationsElement = document.createElement('span');
    courseClassificationsElement.classList.add('course-classifications');
    courseClassificationsElement.innerHTML = course.classification_icons || "";

    const courseLanguagesElement = document.createElement('span');
    courseLanguagesElement.classList.add('course-languages');
    courseLanguagesElement.innerHTML = course.lang_icons || "";

    courseIconsElement.appendChild(courseClassificationsElement);
    courseIconsElement.appendChild(courseLanguagesElement);

    courseElement.appendChild(courseIconsElement);

    const courseGroupElement = document.querySelector(`.course-group[data-group="${group}"]`);
    courseGroupElement.appendChild(courseElement);
}

/**
 * Helper function to populate filters (either classifications or languages).
 * 
 * @param {Array} items - Array of either classifications or languages.
 * @param {string} groupClass - The class of the filter group ('class' or 'lang').
 */
const populateFiltersHelper = (items, groupClass) => {
    const groupContainer = document.querySelector(`.filter-group--${groupClass}`);
    for (const item of items) {
        const filterItem = document.createElement('div');
        filterItem.classList.add('filter-item', `filter-item--${item.abbreviation}`);
        filterItem.setAttribute('data-classification', `${item.abbreviation}`);

        const iconWrapper = document.createElement('span');
        iconWrapper.classList.add('icon-wrapper');
        iconWrapper.innerHTML = item.icon_html; // Ensure this content is sanitized/trusted!

        const filterText = document.createElement('span');
        filterText.classList.add('filter-text');
        filterText.textContent = item.longname;

        filterItem.appendChild(iconWrapper);
        filterItem.appendChild(filterText);

        groupContainer.appendChild(filterItem);
    }
}

/**
 * Populates the classifications filters.
 * 
 * @param {CoursesData} data - The courses data containing classifications.
 */
const populateFilters = (data) => {
    populateFiltersHelper(data.classifications, 'class');
    populateFiltersHelper(data.languages, 'lang');
}

/**
 * Sets the default filters to be selected on page load.
 */
const setDefaultFilters = () => {
    toggleFilter(document.querySelector(".filter-item--imp"), "imp");
    toggleFilter(document.querySelector(".filter-item--theo"), "theo");

    document.querySelector("#main-course-list").setAttribute("data-loading", "false");
    document.querySelector(".courses-menu").setAttribute("data-loading", "false");
}