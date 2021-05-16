const puppeteer = require('puppeteer');

(async () => {
  const browser = await puppeteer.launch({
    defaultViewport: {
      width: 960, height: 540
    }
  });
  const page = await browser.newPage();

  try {
    await page.goto('http://localhost:3000/');
    await page.screenshot({ path: '../public/dashboard.png' });
  } catch (err) {
    // エラーが起きた際の処理
  } finally {
    await browser.close();
  }
})();