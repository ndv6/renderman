import 'package:puppeteer/puppeteer.dart';
import 'package:dotenv/dotenv.dart' show env;
import 'package:shelf/shelf.dart' as shelf;
import 'utils.dart';

var skippedResources = [
  'googletagmanager',
  'google-analytics',
  'facebook',
  'twitter',
  'windows',
  'service-worker'
];

Future<Content> fetchHeadless(String uri, shelf.Request request) async {
  Browser browser;
  Page page;

  if (env['CHROMIUM_URL'] != null) {
    browser = await puppeteer.connect(browserWsEndpoint: env['CHROMIUM_URL']);
  } else {
    browser = await puppeteer.launch(
      // headless: false,
      // devTools: true,
      ignoreHttpsErrors: true,
      noSandboxFlag: true,
      userDataDir: './tmp/chrome',
      args: [
        '--disable-web-security',
        '--no-sandbox',
        '--disable-setuid-sandbox',
        '--disable-dev-shm-usage',
        '--disable-accelerated-2d-canvas',
        '--disable-gpu',
        '--disable-notifications',
      ],
    );
  }

  page = await browser.newPage();
  await page.setRequestInterception(true);
  page.onRequest.listen((Request request) {
    for (String item in skippedResources) {
      if (request.url.contains(item)) {
        request.abort(error: ErrorReason.blockedByClient);
        break;
      }
    }

    if (request.resourceType == ResourceType.image) {
      request.abort(error: ErrorReason.blockedByClient);
    } else {
      request.continueRequest(headers: request.headers);
    }
  });

  await page.goto(uri, wait: Until.networkIdle);
  // await page.close();

  Content resp = new Content();
  resp.content = await page.content;
  resp.headers = await {'content-type': 'text/html'};

  await browser.disconnect();

  return resp;
}
