import { test, expect } from '@playwright/test';
import { format, addMonths, subMonths, getDaysInMonth } from 'date-fns';

const BASE_URL = 'http://localhost:34115/'; // Assuming this is the correct base URL

test.describe('Progress Page Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(BASE_URL);
    // Navigate to the Progress page. 
    // Assuming there's a link or button with text "Progress" to click.
    // Adjust the selector if needed.
    await page.getByRole('link', { name: 'Progress' }).click();
  });

  test('should display the current month and year initially', async ({ page }) => {
    const currentDate = new Date();
    const expectedMonthYear = format(currentDate, 'MMMM yyyy');
    // Assuming the month/year heading is an <h3> element.
    // The previous changes put it inside a div: <div class="flex flex-row align-center justify-between"><h3>...</h3>...</div>
    // We'll target the h3 that contains the month and year.
    await expect(page.locator('h3', { hasText: expectedMonthYear })).toBeVisible();
  });

  test('should navigate to the next month when "Next" button is clicked', async ({ page }) => {
    const initialDate = new Date();
    const nextMonthDate = addMonths(initialDate, 1);
    const expectedMonthYear = format(nextMonthDate, 'MMMM yyyy');

    // Click the "Next" button
    await page.getByRole('button', { name: 'Next +' }).click();
    
    // Check if the month/year heading updated
    await expect(page.locator('h3', { hasText: expectedMonthYear })).toBeVisible();

    // Check if the number of days displayed matches the next month
    // This relies on the fact that each day <rect> has a child <text> element with the day number.
    const daysInNextMonth = getDaysInMonth(nextMonthDate);
    // Assuming the day texts are direct children of 'g.boxes' or similar structure
    // The current structure is <g class="boxes"> ... <text>DayNumber</text> ... </g>
    // We look for text elements that are numbers from 1 to daysInNextMonth.
    // This is a proxy for checking rendered day elements.
    const dayElements = page.locator('g.boxes text');
    await expect(dayElements).toHaveCount(daysInNextMonth);
  });

  test('should navigate to the previous month when "Prev" button is clicked', async ({ page }) => {
    const initialDate = new Date();
    const previousMonthDate = subMonths(initialDate, 1);
    const expectedMonthYear = format(previousMonthDate, 'MMMM yyyy');

    // Click the "Prev" button
    await page.getByRole('button', { name: '- Prev' }).click();
    
    // Check if the month/year heading updated
    await expect(page.locator('h3', { hasText: expectedMonthYear })).toBeVisible();

    // Check if the number of days displayed matches the previous month
    const daysInPreviousMonth = getDaysInMonth(previousMonthDate);
    const dayElements = page.locator('g.boxes text');
    await expect(dayElements).toHaveCount(daysInPreviousMonth);
  });

  test('navigation buttons should update graph content (check day count)', async ({ page }) => {
    let currentDate = new Date();
    
    // Go to Next Month
    await page.getByRole('button', { name: 'Next +' }).click();
    currentDate = addMonths(currentDate, 1);
    let expectedDays = getDaysInMonth(currentDate);
    await expect(page.locator('g.boxes text')).toHaveCount(expectedDays);
    await expect(page.locator('h3', { hasText: format(currentDate, 'MMMM yyyy') })).toBeVisible();

    // Go to Next Month again
    await page.getByRole('button', { name: 'Next +' }).click();
    currentDate = addMonths(currentDate, 1);
    expectedDays = getDaysInMonth(currentDate);
    await expect(page.locator('g.boxes text')).toHaveCount(expectedDays);
    await expect(page.locator('h3', { hasText: format(currentDate, 'MMMM yyyy') })).toBeVisible();

    // Go to Previous Month
    await page.getByRole('button', { name: '- Prev' }).click();
    currentDate = subMonths(currentDate, 1);
    expectedDays = getDaysInMonth(currentDate);
    await expect(page.locator('g.boxes text')).toHaveCount(expectedDays);
    await expect(page.locator('h3', { hasText: format(currentDate, 'MMMM yyyy') })).toBeVisible();
  });

});
