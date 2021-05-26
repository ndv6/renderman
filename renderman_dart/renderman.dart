import 'package:puppeteer/puppeteer.dart';
import 'package:dotenv/dotenv.dart' show env;
import 'package:shelf/shelf.dart' as shelf;
import 'package:pedantic/pedantic.dart';
import 'package:http/http.dart' as http;
import 'utils.dart';

const skippedResources = [
  'googletagmanager',
  'google-analytics',
  'facebook',
  'twitter',
  'windows',
  'service-worker'
];

const skippedResourceType = [ResourceType.image, ResourceType.font];

Browser browser;

Future<Browser> getBrowserInstance() async {
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

  return browser;
}

Future<Content> fetchHeadless(String uri, shelf.Request request) async {
  Page page;

  if (browser == null) {
    browser = await getBrowserInstance();
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

    if (skippedResourceType.contains(request.resourceType)) {
      request.abort(error: ErrorReason.blockedByClient);
    } else {
      request.continueRequest(headers: request.headers);
    }
  });

  var response = await page.goto(uri, wait: Until.networkIdle);
  String body = await page.content;

  Content resp = new Content(response.headers, body);

  await page.close();

  return resp;
}

Future<Content> fetchHttp(String uri, shelf.Request request) async {
  var httpClientRequest = http.StreamedRequest('GET', Uri.parse(uri));

  unawaited(request
      .read()
      .forEach(httpClientRequest.sink.add)
      .catchError(httpClientRequest.sink.addError)
      .whenComplete(httpClientRequest.sink.close));

  var response = await http.Client().send(httpClientRequest);

  response.headers.remove('transfer-encoding');

  if (response.headers['content-encoding'] == 'gzip') {
    response.headers.remove('content-encoding');
    response.headers.remove('content-length');
  }

  return new Content(response.headers, response.stream);
}
