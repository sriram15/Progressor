import { test, expect } from '@playwright/test';

const BASE_URL = 'http://localhost:34115/'
test('can create a new card', async ({ page }) => {
  await page.goto(BASE_URL);

  // Expect a title "to contain" a substring.
  await expect(page).toHaveTitle(/wails-creator-todo/);

  await expect(page.getByText('Progressor')).toBeVisible();

  await expect(page.getByRole('link', { name: 'Home' })).toBeVisible();
  await expect(page.getByRole('button', { name: 'Add Card' })).toBeVisible();


  await page.getByRole('button', { name: 'Add Card' }).click();
  await page.getByRole('heading', { name: 'Add Card Title' }).click();
 
  await page.getByRole('dialog').click();
  await page.getByPlaceholder('Enter card title').fill('testing');
  
  await page.getByRole('button', { name: 'Submit' }).click();
  await page.getByRole('heading', { name: 'testing' }).click();
});
