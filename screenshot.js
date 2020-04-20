const puppeteer = require('puppeteer');

let args = process.argv.slice(2);

(async () => {

  // Use Pi's chromium-browser, rather than Puppeteer's chrome which doesn't work on the Pi's ARM architecture: https://github.com/puppeteer/puppeteer/issues/4249#issuecomment-535727445
  const browser = await puppeteer.launch({
    // headless: false,
    executablePath: 'chromium-browser'
    });
  const page = await browser.newPage();

  process.stdout.write('Loading page...');
  await page.goto('https://magicseaweed.com/' + args[0], {waitUntil: 'networkidle2'});
  process.stdout.write('done.\n');

  // Fix blank screenshots for elements outside of viewport: https://github.com/puppeteer/puppeteer/issues/2423#issuecomment-590738707
  const viewport = { width: 1440, height: 900, deviceScaleFactor: 2 };
  const fullPage = await page.$('body');
  const fullPageSize = await fullPage.boundingBox();
  await page.setViewport(
    Object.assign({}, viewport, { height: fullPageSize.height })
  );

  // Hide.
  // await page.$eval('#msw-js-fc > div.table-responsive-xs > table > tbody:nth-child(2) > tr.background-clear.msw-js-tide', e => e.setAttribute("style", "display: none"));
  // await page.$eval('#msw-js-fct > header > div', e => e.setAttribute("style", "display: none"));
  await page.$eval('#corona-message-container', e => e.setAttribute("style", "display: none"));

  let elements = {
  // 'dayTitle': 'body > div.cover > div.cover-inner > div.pages.clear-left.clear-right > div > div.msw-fc.msw-js-forecast > div:nth-child(2) > div:nth-child(2) > div > div > div.msw-col-fluid > div > header > h3',
  'current': 'body > div.cover > div.cover-inner > div.pages.clear-left.clear-right > div > div.msw-fc.msw-js-forecast > div:nth-child(2) > div:nth-child(2) > div > div > div.msw-col-fluid > div > div.row.margin-bottom',
  // 'weekTitle': '#msw-js-fct > header',
  'weekForecast': '#tab-7day > div',
  // 'header': '#msw-js-fc > div.table-responsive-xs > table > thead',
  // 'dayForecast': '#msw-js-fc > div.table-responsive-xs > table > tbody:nth-child(2)'
  };

  process.stdout.write('Screenshitting...');
  for(var name in elements) {
    var selector = elements[name];
    elementHandle = await page.$(selector);
    await elementHandle.screenshot({path: '/home/pi/spindrift/img/' + name + '.png'});
  }
  process.stdout.write('done.\n');

  browser.close();
})();
