const puppeteer = require('puppeteer');

let args = process.argv.slice(2);

(async () => {

  // Use Pi's chromium-browser, rather than Puppeteer's chrome which doesn't work on the Pi's ARM architecture: https://github.com/puppeteer/puppeteer/issues/4249#issuecomment-535727445
  const browser = await puppeteer.launch({
    // headless: false,
    executablePath: 'chromium-browser'
    });
  const page = await browser.newPage();

  process.stdout.write('Loading surf forecast...');
  try {
    await page.goto('https://magicseaweed.com/' + args[0], {waitUntil: 'networkidle2'});
    process.stdout.write('done.\n');
  }
  catch(err) {
    process.stderr.write(err.message);
    process.exit(1);
  }

  // Fix blank screenshots for elements outside of viewport: https://github.com/puppeteer/puppeteer/issues/2423#issuecomment-590738707
  const viewport = { width: 1440, height: 900, deviceScaleFactor: 2 };
  const fullPage = await page.$('body');
  const fullPageSize = await fullPage.boundingBox();
  await page.setViewport(
    Object.assign({}, viewport, { height: fullPageSize.height })
  );

  // Prepare page by hiding/resizing elements.
  await page.$eval('#corona-message-container', e => e.setAttribute('style', 'display: none'));
  for (let i = 1; i <= 7; i++) {
    await page.$eval('.scrubber-forecast-graph-wind > div:nth-child(2) > div:nth-child(1) > div:nth-child(' + i + ') > div:nth-child(1) > div:nth-child(1)', e => e.setAttribute('style', 'display: none'));
    await page.$eval('.scrubber-forecast-graph > div:nth-child(2) > div:nth-child(1) > div:nth-child(' + i + ') > div:nth-child(2)', e => e.setAttribute('style', 'padding-top: 0px'));
  }

  let elements = {
    'current': 'body > div.cover > div.cover-inner > div.pages.clear-left.clear-right > div > div.msw-fc.msw-js-forecast > div:nth-child(2) > div:nth-child(2) > div > div > div.msw-col-fluid > div > div.row.margin-bottom',
    'surf-forecast': '#tab-7day > div',
  };

  process.stdout.write('Screenshitting:\n');
  for(var name in elements) {
    var selector = elements[name];
    elementHandle = await page.$(selector);
    process.stdout.write('- ' + name + '\n');
    await elementHandle.screenshot({path: '/home/pi/spindrift/img/' + name + '.png'});
  }

  browser.close();
})();
