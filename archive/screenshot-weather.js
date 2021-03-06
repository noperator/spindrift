const puppeteer = require('puppeteer');

let args = process.argv.slice(2);

(async () => {

  // Use Pi's chromium-browser, rather than Puppeteer's chrome which doesn't work on the Pi's ARM architecture: https://github.com/puppeteer/puppeteer/issues/4249#issuecomment-535727445
  const browser = await puppeteer.launch({
    // headless: false,
    executablePath: 'chromium-browser'
    });
  const page = await browser.newPage();

  process.stdout.write('Loading weather forecast...');
  try {
    await page.goto('https://www.google.com/search?q=weather forecast ' + args[0], {waitUntil: 'networkidle2'});
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
  await page.$eval('div.wob_df:nth-child(8)', e => e.setAttribute('style', 'display: none'));
  await page.$eval('div.gic:nth-child(5)', e => e.setAttribute('style', 'width: calc(75px * 7); height: 72px'));
  for (let i = 1; i <= 7; i++) {
    await page.$eval('div.wob_df:nth-child(' + i + ')', e => e.setAttribute('class', 'wob_df'));
    await page.$eval('div.wob_df:nth-child(' + i + ') > div:nth-child(1)', e => e.setAttribute('style', 'display: none'));
  }

  let elements = {
    'temp': '#wob_gsp',
    'rain': '#wob_gsp',
    'wind': '#wob_gsp',
    'weather-forecast': 'div.gic:nth-child(5)'
  };

  process.stdout.write('Screenshitting:\n');
  for (let name in elements) {
    let selector = elements[name];
    elementHandle = await page.$(selector);
    switch (name) {
      case 'temp':
        await page.click('#wob_temp');
        break;
      case 'rain':
        await page.click('#wob_rain');
        break;
      case 'wind':
        await page.click('#wob_wind');
        break;
    }
    process.stdout.write('- ' + name + '\n');
    await elementHandle.screenshot({path: '/home/pi/spindrift/img/' + name + '.png'});
  }

  browser.close();
})();
