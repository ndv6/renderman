class Content {
  Map<String, String> headers;
  dynamic content;
}

bool isAssetRequest(String uri) {
  RegExp exp = new RegExp(
      r"(.*(css|js|json|jpg|jpeg|gif|svg|png|css|mp4|ogg|svg|ttf|ico|map))");

  return exp.hasMatch(uri);
}

bool isBotRequest(String agent) {
  RegExp exp = new RegExp(
      r"(googlebot|bingbot|yandex|twitterbot|facebookexternalhit|linkedinbot|embedly|pinterest\/0\.|pinterestbot|slackbot|whatsapp|curl|Slack-ImgProxy|Slack)");

  return exp.hasMatch(agent);
}
