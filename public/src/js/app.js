/*
    app.js
    Authored by Ryan Pope <ryan@rcpope.net>
    Licensed under the GNU Public License, v3.0
*/

/**
 * Changes a CSS property, for a set of elements, for a specified period of time
 * 
 * @param {string} selector The selector matching the element(s).
 * @param {string} property The property to set.
 * @param {string} value The value to temporarily set the property to.
 * @param {number} delay The delay between setting the new property and resetting the old one.
 */
const temporaryStyle = (selector, property, value, delay) => {
    document.querySelectorAll(selector).forEach((element) => {
        const computedValue = window.getComputedStyle(element).getPropertyValue(property);
        element.style.setProperty(property, value);
        setTimeout(() => {
            if(value != computedValue){
                element.style.setProperty(property, computedValue);
            }
        }, delay);
    });
}

/**
 * Simulates a fade-in effect by smoothly changing an element's opacity from 0 to 1.
 *
 * @param {HTMLElement} element - The DOM element to be faded in.
 * @param {number} [duration=1000] - The duration of the fade-in effect in milliseconds.
 *
 * @example
 * fadeIn(document.getElementById('myElement'), 1500);
 *
 */
function fadeIn(element, duration = 1000) {
    element.style.display = 'block';
    element.style.opacity = 0;

    let last = +new Date();
    let tick = function() {
        element.style.opacity = +element.style.opacity + (new Date() - last) / duration;
        last = +new Date();

        if (+element.style.opacity < 1) {
            (window.requestAnimationFrame && requestAnimationFrame(tick)) || setTimeout(tick, 16);
        }
    };

    tick();
}