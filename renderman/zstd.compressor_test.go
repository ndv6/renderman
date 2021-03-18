package main

import (
	"testing"
	"time"
)

var zstdc = newZstdCompressor()

var html = `<html lang="id"><head><style type="text/css">
.anticon {
  display: inline-block;
  color: inherit;
  font-style: normal;
  line-height: 0;
  text-align: center;
  text-transform: none;
  vertical-align: -0.125em;
  text-rendering: optimizeLegibility;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.anticon > * {
  line-height: 1;
}

.anticon svg {
  display: inline-block;
}

.anticon::before {
  display: none;
}

.anticon .anticon-icon {
  display: block;
}

.anticon[tabindex] {
  cursor: pointer;
}

.anticon-spin::before,
.anticon-spin {
  display: inline-block;
  -webkit-animation: loadingCircle 1s infinite linear;
  animation: loadingCircle 1s infinite linear;
}

@-webkit-keyframes loadingCircle {
  100% {
    -webkit-transform: rotate(360deg);
    transform: rotate(360deg);
  }
}

@keyframes loadingCircle {
  100% {
    -webkit-transform: rotate(360deg);
    transform: rotate(360deg);
  }
}
</style><meta charset="utf-8"><link rel="icon" href="/favicon-prospeku.ico"><meta name="viewport" content="width=device-width,initial-scale=1"><meta name="theme-color" content="#000000"><link rel="apple-touch-icon" href="/logo192.png"><link rel="manifest" href="/manifest.json"><title>Prospeku</title><script async="" src="https://www.googletagmanager.com/gtm.js?id=GTM-KL83PS6"></script><script>!function(e,t,a,n,g){e[n]=e[n]||[],e[n].push({"gtm.start":(new Date).getTime(),event:"gtm.js"});var m=t.getElementsByTagName(a)[0],r=t.createElement(a);r.async=!0,r.src="https://www.googletagmanager.com/gtm.js?id=GTM-KL83PS6",m.parentNode.insertBefore(r,m)}(window,document,"script","dataLayer")</script><link href="/static/css/2.1a0e7f37.chunk.css" rel="stylesheet"><link href="/static/css/main.63c072ca.chunk.css" rel="stylesheet"><style type="text/css">.TransformComponent-module_container__3NwNd {
  position: relative;
  width: fit-content;
  height: fit-content;
  overflow: hidden;
  -webkit-touch-callout: none; /* iOS Safari */
  -webkit-user-select: none; /* Safari */
  -khtml-user-select: none; /* Konqueror HTML */
  -moz-user-select: none; /* Firefox */
  -ms-user-select: none; /* Internet Explorer/Edge */
  user-select: none;
  margin: 0;
  padding: 0;
}
.TransformComponent-module_content__TZU5O {
  display: flex;
  flex-wrap: wrap;
  width: fit-content;
  height: fit-content;
  margin: 0;
  padding: 0;
  transform-origin: 0% 0%;
}
.TransformComponent-module_content__TZU5O img {
  pointer-events: none;
}
</style></head><body><noscript>You need to enable JavaScript to run this app.</noscript><noscript><iframe src="https://www.googletagmanager.com/ns.html?id=GTM-KL83PS6" height="0" width="0" style="display:none;visibility:hidden"></iframe></noscript><div id="root"><div class="App"><div class="wrap"><section class="ant-layout"><header class="ant-layout-header"><div class="header" id="stickyHeader"><div type="flex" class="ant-row ant-row-middle header-style" style="row-gap: 0px;"><div class="ant-col ant-col-xs-22 ant-col-xs-offset-1 ant-col-sm-22 ant-col-sm-offset-1 ant-col-md-20 ant-col-md-offset-2 ant-col-lg-20 ant-col-lg-offset-2 ant-col-xl-18 ant-col-xl-offset-3"><div type="flex" class="ant-row ant-row-middle" style="margin-left: -8px; margin-right: -8px; row-gap: 0px;"><div class="ant-col d-lg-none d-block text-center" style="padding-left: 8px; padding-right: 8px;"><div><button type="button" class="ant-btn ant-btn-link menu-trigger"><svg width="25" height="25" viewBox="0 0 21 15" type="menu"><defs><path id="prefix__a" d="M0 1.116L20 1.116 20 3.616 0 3.616z"></path></defs><g fill="none" fill-rule="evenodd"><path fill="#40A280" d="M18.75 8.817H1.25C.56 8.817 0 8.258 0 7.567s.56-1.25 1.25-1.25h17.5c.69 0 1.25.559 1.25 1.25s-.56 1.25-1.25 1.25" transform="translate(.75 -.5)"></path><g transform="translate(.75 -.5) translate(0 .134)"><mask id="prefix__b" fill="#fff"><use xlink:href="#prefix__a"></use></mask><path fill="#40A280" d="M18.75 3.616H1.25C.56 3.616 0 3.057 0 2.366c0-.692.56-1.25 1.25-1.25h17.5c.69 0 1.25.558 1.25 1.25 0 .691-.56 1.25-1.25 1.25" mask="url(#prefix__b)"></path></g><path fill="#40A280" d="M18.75 13.884H1.25c-.69 0-1.25-.559-1.25-1.25s.56-1.25 1.25-1.25h17.5c.69 0 1.25.559 1.25 1.25s-.56 1.25-1.25 1.25" transform="translate(.75 -.5)"></path></g></svg></button></div></div><div class="ant-col ant-col-xs-18 ant-col-sm-18 ant-col-md-5 ant-col-lg-5 ant-col-xl-5" style="padding-left: 8px; padding-right: 8px;"><a class="header__logo text-center" href="/"><svg viewBox="0 0 207.28 39.37" height="32" width="194"><defs><style>.cls-1{fill:#efdd00;}.cls-2{fill:#f4eb9f;}.cls-3{fill:#9dbda5;}.cls-4{fill:#cedbdd;}.cls-5{fill:#246755;}.cls-6{fill:#8bcc96;}.cls-7{fill:#3db75a;}</style></defs><title>Asset 22</title><g id="Layer_2" data-name="Layer 2"><g id="Layer_1-2" data-name="Layer 1"><path class="cls-1" d="M19.86,25.41l-.08-.09.06-.06S19.85,25.36,19.86,25.41Z"></path><path class="cls-1" d="M45.64,15a15,15,0,0,1-5.74,11.8,2.19,2.19,0,0,0-.52-.84l-8.6-8.59a2.15,2.15,0,0,0-2.07-.56,2.08,2.08,0,0,0-1,.56l-6,6v-19a.17.17,0,0,0,0-.07.68.68,0,0,0,0-.19.14.14,0,0,0,0-.06,1.54,1.54,0,0,0-.05-.17l0-.06-.05-.11a.05.05,0,0,1,0,0l-.07-.1a.64.64,0,0,0-.1-.13l-.1-.12L21.68,3a15,15,0,0,1,24,12Z"></path><path class="cls-2" d="M7.7,17.54V33.63a22.86,22.86,0,0,0-2.34,1A9.49,9.49,0,0,1,7.7,17.54Z"></path><path class="cls-2" d="M19.86,25.41l-.08-.09.06-.06S19.85,25.36,19.86,25.41Z"></path><path class="cls-3" d="M28.71,6.07V16.8a2.08,2.08,0,0,0-1,.56l-6,6v-19a.17.17,0,0,0,0-.07.68.68,0,0,0,0-.19.14.14,0,0,0,0-.06,1.54,1.54,0,0,0-.05-.17l0-.06-.05-.11a.05.05,0,0,1,0,0l-.07-.1a.64.64,0,0,0-.1-.13l-.1-.12L21.2,3.2l-.06-.05-.06,0A.24.24,0,0,0,21,3L20.86,3a.61.61,0,0,0-.14-.07h0l-.15-.07h0l-.07,0-.12,0A1.35,1.35,0,0,0,20,2.67a1.14,1.14,0,0,1,.45,0L21.68,3l5.85,1.36A1.71,1.71,0,0,1,28.71,6.07Z"></path><path class="cls-4" d="M19.78,25.32l.06-.06,1.93-1.93v-19a.17.17,0,0,0,0-.07.68.68,0,0,0,0-.19.14.14,0,0,0,0-.06,1.54,1.54,0,0,0-.05-.17l0-.06-.05-.11a.05.05,0,0,1,0,0l-.07-.1a.64.64,0,0,0-.1-.13l-.1-.12L21.2,3.2l-.06-.05-.06,0A.24.24,0,0,0,21,3L20.86,3a.61.61,0,0,0-.14-.07h0l-.15-.07h0l-.07,0-.12,0A1.35,1.35,0,0,0,20,2.67a2.23,2.23,0,0,0-.63,0l-10,1.65A1.85,1.85,0,0,0,7.7,6.07V33.63a18.35,18.35,0,0,1,10-.85,18.88,18.88,0,0,1,3.36,1V29h-1.9a2.15,2.15,0,0,1,0-3Zm-1.07-11-3.05.34-4.7.54V12.6a.72.72,0,0,1,.53-.69l4.63-.64L17.84,11a.71.71,0,0,1,.87.69Zm0-4.58L16.5,10,11,10.64V8a.71.71,0,0,1,.53-.69l6.35-.88a.68.68,0,0,1,.45,0,.7.7,0,0,1,.42.65Z"></path><path class="cls-4" d="M19.86,25.41l-.08-.09.06-.06S19.85,25.36,19.86,25.41Z"></path><path class="cls-5" d="M39.38,29H37.47v5.37a1.88,1.88,0,0,1,0,.34,16.84,16.84,0,0,0-6.08-.25V32.74a2.27,2.27,0,0,0-2.28-2.29,2.34,2.34,0,0,0-.38,0,2.3,2.3,0,0,0-1.91,2.26v3c-.46.19-.91.39-1.33.59a22.7,22.7,0,0,0-3.7-2.16c-.25-.12-.5-.22-.75-.33V29h-1.9a2.15,2.15,0,0,1,0-3l.66-.66.08.09c0-.05,0-.1,0-.15l1.93-1.93,6-6a2.08,2.08,0,0,1,1-.56,2.15,2.15,0,0,1,2.07.56L39.38,26a2.19,2.19,0,0,1,.52.84h0A2.14,2.14,0,0,1,39.38,29Z"></path><path class="cls-5" d="M19.86,25.41l-.08-.09.06-.06S19.85,25.36,19.86,25.41Z"></path><path class="cls-4" d="M31.37,32.74v1.72a18.13,18.13,0,0,0-2.66.58,19.42,19.42,0,0,0-1.91.66v-3a2.3,2.3,0,0,1,1.91-2.26,2.34,2.34,0,0,1,.38,0,2.27,2.27,0,0,1,2.28,2.29Z"></path><path class="cls-6" d="M44.39,39.37H29.45c-.77-.76-1.54-1.44-2.3-2.06q-.52-.44-1-.81l-.34-.25c.43-.2.88-.4,1.35-.59A18.94,18.94,0,0,1,29,35a18.54,18.54,0,0,1,2.69-.58,16.79,16.79,0,0,1,6.15.25,19.53,19.53,0,0,1,7,3.21A.83.83,0,0,1,44.39,39.37Z"></path><path class="cls-7" d="M29.45,39.37H.83A.82.82,0,0,1,.3,37.92a26.82,26.82,0,0,1,5.12-3.35,24.26,24.26,0,0,1,2.36-1,18.51,18.51,0,0,1,10.08-.85,18.75,18.75,0,0,1,3.4,1c.25.1.51.21.76.33a22.45,22.45,0,0,1,3.74,2.18l.34.25q.52.38,1,.81C27.91,37.93,28.68,38.61,29.45,39.37Z"></path><path class="cls-3" d="M18.71,7.14V9.76L16.5,10,11,10.64V8a.71.71,0,0,1,.53-.69l6.35-.88a.68.68,0,0,1,.45,0A.7.7,0,0,1,18.71,7.14Z"></path><path class="cls-3" d="M18.71,11.72v2.62l-3.05.34-4.7.54V12.6a.72.72,0,0,1,.53-.69l4.63-.64L17.84,11A.71.71,0,0,1,18.71,11.72Z"></path><path class="cls-5" d="M70.43,10.9a7.63,7.63,0,0,1,4.7,7.36,8,8,0,0,1-1.22,4.43,7.82,7.82,0,0,1-3.48,2.91,13.07,13.07,0,0,1-5.33,1H60.45V33H55.1V9.87h10A12.89,12.89,0,0,1,70.43,10.9Zm-2,10.31a3.61,3.61,0,0,0,1.26-2.95,3.65,3.65,0,0,0-1.26-3,5.6,5.6,0,0,0-3.66-1H60.45v8H64.8A5.6,5.6,0,0,0,68.46,21.21Z"></path><path class="cls-5" d="M85.31,15.62A9.23,9.23,0,0,1,88.92,15v4.75c-.57-.05-1-.07-1.15-.07a4.71,4.71,0,0,0-3.47,1.24,5,5,0,0,0-1.25,3.71V33H77.9V15.22h4.92v2.34A5.7,5.7,0,0,1,85.31,15.62Z"></path><path class="cls-5" d="M94.58,32.07a8.71,8.71,0,0,1-4.73-8,8.91,8.91,0,0,1,1.25-4.72,8.73,8.73,0,0,1,3.48-3.25,11.38,11.38,0,0,1,10,0,8.67,8.67,0,0,1,3.46,3.25,8.91,8.91,0,0,1,1.26,4.72,8.94,8.94,0,0,1-1.26,4.72,8.81,8.81,0,0,1-3.46,3.25,11.38,11.38,0,0,1-10,0Zm8.28-4.39a5,5,0,0,0,1.27-3.58,5,5,0,0,0-1.27-3.58,4.28,4.28,0,0,0-3.25-1.34,4.33,4.33,0,0,0-3.26,1.34,4.93,4.93,0,0,0-1.29,3.58,4.94,4.94,0,0,0,1.29,3.58A4.33,4.33,0,0,0,99.61,29,4.27,4.27,0,0,0,102.86,27.68Z"></path><path class="cls-5" d="M113.77,32.69a11.27,11.27,0,0,1-3.37-1.37l1.72-3.69A10.69,10.69,0,0,0,115,28.86a12,12,0,0,0,3.3.48c2.18,0,3.27-.54,3.27-1.61a1.15,1.15,0,0,0-.89-1.09,13.08,13.08,0,0,0-2.74-.56,27.63,27.63,0,0,1-3.6-.76,5.66,5.66,0,0,1-2.46-1.52,4.33,4.33,0,0,1-1-3.1,4.87,4.87,0,0,1,1-3,6.3,6.3,0,0,1,2.83-2,12.28,12.28,0,0,1,4.41-.72,16.94,16.94,0,0,1,3.74.41,10.09,10.09,0,0,1,3.09,1.14l-1.72,3.66a10.22,10.22,0,0,0-5.11-1.32,5.21,5.21,0,0,0-2.48.46,1.36,1.36,0,0,0-.82,1.19,1.17,1.17,0,0,0,.89,1.15,15.8,15.8,0,0,0,2.84.63,28.74,28.74,0,0,1,3.56.78,5.24,5.24,0,0,1,2.41,1.5,4.26,4.26,0,0,1,1,3,4.7,4.7,0,0,1-1,2.94,6.36,6.36,0,0,1-2.88,2,13,13,0,0,1-4.51.71A17,17,0,0,1,113.77,32.69Z"></path><path class="cls-5" d="M143.32,16.09a8.18,8.18,0,0,1,3.18,3.21,9.64,9.64,0,0,1,1.16,4.8,9.69,9.69,0,0,1-1.16,4.8,8.22,8.22,0,0,1-3.18,3.2,9.05,9.05,0,0,1-4.51,1.14,6.86,6.86,0,0,1-5.34-2.15v8.28h-5.15V15.22h4.92v2A6.81,6.81,0,0,1,138.81,15,9.15,9.15,0,0,1,143.32,16.09Zm-2.15,11.59a5,5,0,0,0,1.27-3.58,5,5,0,0,0-1.27-3.58,4.61,4.61,0,0,0-6.5,0,5,5,0,0,0-1.27,3.58,5,5,0,0,0,1.27,3.58,4.64,4.64,0,0,0,6.5,0Z"></path><path class="cls-5" d="M167.72,25.55H154.29A4.16,4.16,0,0,0,156,28.16a5.72,5.72,0,0,0,3.36.95,6.93,6.93,0,0,0,2.46-.41,6.19,6.19,0,0,0,2-1.3l2.73,3q-2.51,2.86-7.32,2.87a11.59,11.59,0,0,1-5.31-1.17,8.73,8.73,0,0,1-3.57-3.25,9,9,0,0,1-1.25-4.72,9.06,9.06,0,0,1,1.24-4.71,8.6,8.6,0,0,1,3.4-3.26,10.53,10.53,0,0,1,9.55-.05,8.14,8.14,0,0,1,3.31,3.22,9.48,9.48,0,0,1,1.21,4.86C167.82,24.23,167.79,24.69,167.72,25.55Zm-12-5.71a4.2,4.2,0,0,0-1.46,2.71H163a4.26,4.26,0,0,0-1.45-2.69,4.37,4.37,0,0,0-2.91-1A4.42,4.42,0,0,0,155.68,19.84Z"></path><path class="cls-5" d="M176.89,26l-2.47,2.44V33h-5.15V8.49h5.15V22.35l7.52-7.13h6.14l-7.39,7.52L188.74,33H182.5Z"></path><path class="cls-5" d="M207.28,15.22V33H202.4V30.86A6.9,6.9,0,0,1,200,32.63a7.82,7.82,0,0,1-3.07.61,7.52,7.52,0,0,1-5.54-2q-2.06-2-2-6v-10h5.15v9.27q0,4.29,3.59,4.29a3.85,3.85,0,0,0,3-1.2A5.07,5.07,0,0,0,202.14,24V15.22Z"></path></g></g></svg></a></div><div class="ant-col header-pd d-none d-lg-block ant-col-xs-12 ant-col-sm-12 ant-col-md-19 ant-col-lg-19 ant-col-xl-19" style="padding-left: 8px; padding-right: 8px;"><div class="header__menu"><a href="/">Beranda</a><a href="/#" class="ant-dropdown-trigger ant-dropdown-link">Produk</a><a href="/artikel">Artikel</a><a href="/panduan">Panduan</a><a href="https://prospeku-cms-sit.onelabs.dev/register" class="ant-btn ant-btn-md ant-btn-secondary btn-register">Daftarkan Kantor Agen</a><a href="https://cms.prospeku.com/login" class="ant-btn ant-btn-md ant-btn-primary btn-login">Masuk</a></div></div></div></div></div></div><div class="sidebar d-lg-none d-block false"><ul class="sidebar__menu"><li><a href="/">Beranda</a></li><li><div class="ant-collapse ant-collapse-borderless ant-collapse-icon-position-left collapse-custom"><div class="ant-collapse-item site-collapse-custom-panel"><div class="ant-collapse-header" role="button" tabindex="0" aria-expanded="false"><span role="img" aria-label="caret-right" class="anticon anticon-caret-right ant-collapse-arrow"><svg viewBox="0 0 1024 1024" focusable="false" data-icon="caret-right" width="1em" height="1em" fill="currentColor" aria-hidden="true"><path d="M715.8 493.5L335 165.1c-14.2-12.2-35-1.2-35 18.5v656.8c0 19.7 20.8 30.7 35 18.5l380.8-328.4c10.9-9.4 10.9-27.6 0-37z"></path></svg></span>Produk</div></div></div></li><li><a href="/artikel">Artikel</a></li></ul></div></header><main class="ant-layout-content"><div type="flex" class="ant-row ant-row-center ant-row-middle" style="height: 100vh; row-gap: 0px;"><div class="ant-spin ant-spin-lg ant-spin-spinning"><span class="ant-spin-dot ant-spin-dot-spin"><i class="ant-spin-dot-item"></i><i class="ant-spin-dot-item"></i><i class="ant-spin-dot-item"></i><i class="ant-spin-dot-item"></i></span></div></div></main><footer class="ant-layout-footer"><div type="flex" middle="center" class="ant-row ant-footer footer--custom" style="row-gap: 0px;"><div class="ant-col ant-col-xs-22 ant-col-xs-offset-1 ant-col-sm-22 ant-col-sm-offset-1 ant-col-md-20 ant-col-md-offset-2 ant-col-lg-20 ant-col-lg-offset-2 ant-col-xl-20 ant-col-xl-offset-2"><div type="flex" class="ant-row ant-row-middle" style="margin-left: -8px; margin-right: -8px; row-gap: 0px;"><div class="ant-col text-center text-lg-left mb-4 mb-lg-0 ant-col-xs-24 ant-col-sm-24 ant-col-lg-17" style="padding-left: 8px; padding-right: 8px;"><div class="footer--menu mb-0 mb-lg-4"><a href="/faq">FAQ</a><a href="/hubungi-kami">Hubungi Kami</a><a href="/syarat-dan-ketentuan">Syarat dan Ketentuan</a><a href="/kebijakan-privasi">Kebijakan Privasi</a></div></div></div><div type="flex" class="ant-row ant-row-middle mb-2" style="margin-left: -8px; margin-right: -8px; row-gap: 0px;"><div class="ant-col text-center text-lg-left mb-4 mb-lg-0 ant-col-xs-24 ant-col-sm-24 ant-col-lg-12" style="padding-left: 8px; padding-right: 8px;"><a href="https://www.facebook.com/prospeku" rel="noopener noreferrer nofollow external" target="_blank" class="ant-btn ant-btn-link footer--social"><img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAYAAABXAvmHAAAAAXNSR0IArs4c6QAAAERlWElmTU0AKgAAAAgAAYdpAAQAAAABAAAAGgAAAAAAA6ABAAMAAAABAAEAAKACAAQAAAABAAAAMKADAAQAAAABAAAAMAAAAADbN2wMAAACNElEQVRoBWNgGAWjITAaAqMhMGhD4P///zJAvAyInwExvQHITpDdMvgCiBGXJFTjRaC8EC41dBJ/B7RHn5GR8Qk2+5iwCULFuoD0QDse5BSQG0BuwQrwxcAzoA5JrLroL/gcGANS2KzF54H/2DQMlBjQA1jdii8JDZRbSbKXhSTVJCr+8uMbw9VH9xkevnrB8PLjewaG/5BItdTQZTBR0SDRNOzKaeaBDScOMbSvWcTw7ssnDJtZWVgGtwfO3LnBULpgCobDaSFAkzwwectqWrgVq5k0SUI3nz7EsMzXzIZBXUqOgYmJicFUVRNDnlwBqnvg5+9fDG8/o6Z7aWFRhr6kPHLdiFcf1ZPQrz9/MCzk5eTCEKOWANU9gM1hOOogbEpJFsNau4FMATU9STFt+vZ1DNOAGKQLlIyQAcgD7KyscKHSwGiGOEdPOJ8YBq6amGp54M/fvww/fqE6HOYwUFggy4nw8cOkKKbpkoTQXSkhIIwuRDafeh7A3tbC6jAJQep5gGp5AObSz9+/MRgVJsC4YFpTVoFhUzXOJj2KWlwcXHmAejGAy2Yai496gMYBTND40RggGEQ0VjAaAzQOYILGj8YAwSCisYIhHwNUa43CAhrUbM7zCYVxwbQIvwAKn5ocqreFqOk4ZLNGZFvoOXIIDDAbp1vwZeIDA+xoZOtxugVfHgDNjAzdCQ5gpnkC9IA+EC8HYpxRCJSjFQDZCbIb5+wMrSweNXc0BEZDYDQEiA8BAMkY+PkozxokAAAAAElFTkSuQmCC" alt="Facebook"></a><a href="https://www.instagram.com/prospeku_id/" rel="noopener noreferrer nofollow external" target="_blank" class="ant-btn ant-btn-link footer--social"><img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAYAAABXAvmHAAAAAXNSR0IArs4c6QAAAERlWElmTU0AKgAAAAgAAYdpAAQAAAABAAAAGgAAAAAAA6ABAAMAAAABAAEAAKACAAQAAAABAAAAMKADAAQAAAABAAAAMAAAAADbN2wMAAADR0lEQVRoBe1ZS2gTURS9k9o21bpQpFCJVmuoFGt/G0FEdOHGjaK4tLjzt+jKKgqCSEW7VshKioJaUdSFuHEhKkI3ltaWYiRtjG2jpShoP2njZHznJVPSZN7EeTOTKMwNk0nuu+/cc+775A0h8syrgFcBrwL/bAU0TQuw6x67pthVbENO5A6YFUgRNWY6DrL29aKYIvm/szwtiqJMGOXzGTkzvh52LzV5UAEHcDE0sxGYYj1qDXsV3xlnI7DRKK2ZAM2oQ6l8TIAhV7MpVCqulvJ6AiyVy4XgVXYwo9Nxuv74Lr2PfKQfs78sQa2rXkvt27bThaPHaUuN/F5huDDABL9aZoxA/lD3eZpfTJiFFWxbXemnZ5duFBQhWsTSI4DKy5A/c/AINdcFCXvKwPgnCr14wkcxdLqroFijAGkBmDYy1rq1gfbvbOddU5lBlsUCiPQuZHXO62I1yp+ZsljAlB4BnZDV+4dohBT2gpAh9tmuuSZg04Ya2t3YzPm9Gx2iLzPT/PPN54/scl7R37aAan8VB0RF5xLpHenwrr10reMUlZel4ZPqb7p4J0RP+1/z2DV+Px8FfJlNLHCf42/YRs0sePKY1trZsRwSnoxp8B243KmpKXXZr3+AD22IQaxuwICvkIkESi/iXMCR2Dh37WtqI5+SDwsf2mDDsTF+d+ItP5Mkqn5WTKqqEEFvwyJ2yhwTsGNzPef0dnSQMOdzDT60wZrq0rG5MTLfHRMQrA1QW30DRb/F6cqD2ytEgDx8aEMMYp0y27tQNpGeE2f5+ajvzUvqD4/QnsYW3ozKgzzOPYhx0hwVgFMlDmZdvbdoYCzMSetkUXmQt3Py1LGy744KADAI9p27SpGvkzT8Ob3bYM47OW0cEYDz/M/5OUokl8hfXpGNyU6aCidciDT6LiwtErBkTXoR42FETaWo+2EvF2GVAMijLzCAJWvCDRm/jGag2Q80ZT4fVVVUmoXntaHyIG/3gUZaABgV85FS9ERmS0BeWV10iARIrwEXuVqC9gRYKpcLwWYjEHchnyykkIuZgFey2VzoJ+RitgvhyPj//sHBtq0JJgDHyfvsEg4ha3PLkBO5hf/OuJXYw/Uq4FXAq8DfV+APnGMIdzkOG90AAAAASUVORK5CYII=" alt="Instagram"></a></div><div class="ant-col text-center text-lg-right ant-col-xs-24 ant-col-sm-24 ant-col-lg-12" style="padding-left: 8px; padding-right: 8px;"><a href="https://prospeku.onelink.me/dWWj/365b478b" target="_blank" rel="noopener noreferrer nofollow external"><img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAIwAAAAqCAYAAABhhbPNAAAAAXNSR0IArs4c6QAAAERlWElmTU0AKgAAAAgAAYdpAAQAAAABAAAAGgAAAAAAA6ABAAMAAAABAAEAAKACAAQAAAABAAAAjKADAAQAAAABAAAAKgAAAADnAesRAAAOkElEQVR4Ae2cBbBVVRfHl4iKBQaiGISKgTp2gyi2owgf4WBgd2GMhfHpYAwqdgcGWIBF2Ci2jjVKGYCFjQoKiiL7+/+W7DPnXl6c9+59z8fnWTPnnnP32b3XXmvttf73mi1I2yppsK7pukJ+/SvngLWHB+CFSqmV3gzTNVdXzij5HMAD8AI8AW8UUGt9G68rZ5R8DirigXHiDXjEFp1/DdK9Iwk55TNQwQy0UFpbXUMX0UcHXWN1NdKVUz4Dlc3APL3ohIQZqGuDynLl6fkMzJ8BhEszPmboajo/Mb/lM1DVDMxEDS10zNKoUSPbe++97YgjjrD99tvPWrd2e6yqgZb0bvnll7cBAwbYEkssUVI9FF5jjTXskksusUUWYa8udNQ0s93SuHFju/LKK+3II4/0UbJoSy655D8yYvpy2WWX+eSvv/76ds899zgD1VVnGGenTp1s0UXR4KXR4osvbh07dlxYGcYaZx3+OuusY6eeeqr9+eefNm/ePLvjjjts8803t5VXXtkeffRR++uvv7JWVXI+dufMmTPtxhtvtG+//dYef/xxu+GGG2zs2LH2yy+/2Iorrmh9+/a1lVZayZ5++ml7+OGHfbGPPfZYu//++23GjBl24okn+hj++OMPO+qoo+zmm2+2du3aWefOnW3ZZZd1ZrzuuutswoQJFkIw8kVabrnl7OSTT7aWLVva888/bw8++KC/ghkOPPBA23LLLe2zzz4zys+aNcvfkb799tvbDz/8YHPmzIlVFdxbtGjh/ab+J554wkaMGOHv99prL1tttdW8f4z90ksvtR9//LGgbH19ySxhmERoscUWs5tuusklzcsvv+yTxoJtu22VTsGyj4eJW2qppbzed955x2bPnm1rr72279xrr73WF3jQoEF2zDHHWJcuXZyht9pqK9t66609H5O+6aabGhth11139Y1AeTbFG2+8YVOmTLGLLrrI64dhItHuwIEDfR7uvPNOO/jgg61Hjx7+Gim00UYbGe2vu+66iTTeY4897Oijj7YhQ4bY3Llzbemll/ZNF+vkjtRkXmEEJOYpp5xiO++8s2fZcccdnREfe+wx3wQnnHBCumi9PmeWMOlJg2nY0ail66+/3ncvO3nw4MGuKthd9U1IPfrTvHlzW2uttezwww+333//3VjUfffd16XQk08+adtss421adPG7rrrLuvQoYN9//33LiUYH+VffPFFe+655+zzzz+33XbbrUB1IEWRWqhBJAzSjHnYf//9bdiwYfbSSy+5dMK+wu6JttXuu+9u9913n7HBvv76a9thhx28Lfocac011zQkDMxI+r333mvdunXzvvD9oYcesldeecVatWrlUjCWq+97ZglT3DGYBmZBzF9zzTW+M9nNr7/+uj9HiVRcrpzfo5rAkFxhhRVs6tSpvsDs4kgwQbQ9WDCkCoyCERulC2olEhIEYsenNwlpvINpiusnL4S02WeffVydIKViXWnG8IwVfFBvWq3T71gv2eMYSMtSXwVNlCWp1gxD63QePX388cfbbbfdZocddpjvQIxjdiqiOk5aWXqbqqRp06YusjkloXruvvtuF+fYCB9++KGdddZZtt1227mkwcaCkBqoA/qNyvn5559tvfXWs3Hj8Hz/PZ54EmLBolHPGHgmDZXx/vvv2znnnONqGFto+PDhXh4JgR3DHTVHW9Do0aNdCu20007eH5i7mBk//fRT++qrr+yMM85whob5Yr3UyQaF6HvsoyfU8wdm/3+ztInRFU9I6fxMIiL3p59+sltuucWNvWik9erVyw3ATz75xKZNm5YuVvIz9kvbtm0NAxHxHQ1PFgKJgUGOkYkof+SRR5L2vvvuO1cNMM8333zjRi0MADEWxgEDwSRIMOwjiHrfeustlwIY1xtuuKEvLMxIG9C7777rUmvjjTe2Z5991hj3Bx984JLv119/NWyZyZMn26uvvmrvvfeel4kf1D9mzBhXmahNNgAqFEK6wOD0mWfUKMb4P0VYdNVeMhg1pspJ4jToZOL1HHTQQUH2Q5JZBmmQ+grS6dW2k6UveZ7q16sO5yhb49UxDNwh3RpOO+00Z4revXuH3377LWEaHmTwhTPPPDNIvOaMk2GT1uGi13r+M5+S1PlqCTF++eWXu2jnLh5x+6JJkyZedpVVVnEfCb4cTgWI2LQBWW0DeYYGMQOZuC2LhEmLEySJRhdkxySS5sILL/Q0nSSCjNMg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" alt="GooglePlay"></a></div></div></div></div><div class="ant-row ant-footer footer--end" style="row-gap: 0px;"><div class="ant-col ant-col-xs-22 ant-col-xs-offset-1 ant-col-sm-22 ant-col-sm-offset-1 ant-col-md-20 ant-col-md-offset-2 ant-col-lg-20 ant-col-lg-offset-2 ant-col-xl-20 ant-col-xl-offset-2"><div type="flex" class="ant-row ant-row-middle" style="margin-left: -8px; margin-right: -8px; row-gap: 0px;"><div class="ant-col ant-col-24" style="padding-left: 8px; padding-right: 8px;">©2020 Prospeku . All rights reserved</div></div></div></div><div class="ant-back-top"></div></footer></section></div></div></div><script>!function(e){function r(r){for(var n,p,l=r[0],a=r[1],f=r[2],c=0,s=[];c<l.length;c++)p=l[c],Object.prototype.hasOwnProperty.call(o,p)&&o[p]&&s.push(o[p][0]),o[p]=0;for(n in a)Object.prototype.hasOwnProperty.call(a,n)&&(e[n]=a[n]);for(i&* Connection #0 to host 777873f5df79.ngrok.io left intact
&i(r);s.length;)s.shift()();return u.push.apply(u,f||[]),t()}function t(){for(var e,r=0;r<u.length;r++){for(var t=u[r],n=!0,l=1;l<t.length;l++){var a=t[l];0!==o[a]&&(n=!1)}n&&(u.splice(r--,1),e=p(p.s=t[0]))}return e}var n={},o={1:0},u=[];function p(r){if(n[r])return n[r].exports;var t=n[r]={i:r,l:!1,exports:{}};return e[r].call(t.exports,t,t.exports,p),t.l=!0,t.exports}p.m=e,p.c=n,p.d=function(e,r,t){p.o(e,r)||Object.defineProperty(e,r,{enumerable:!0,get:t})},p.r=function(e){"undefined"!=typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},p.t=function(e,r){if(1&r&&(e=p(e)),8&r)return e;if(4&r&&"object"==typeof e&&e&&e.__esModule)return e;var t=Object.create(null);if(p.r(t),Object.defineProperty(t,"default",{enumerable:!0,value:e}),2&r&&"string"!=typeof e)for(var n in e)p.d(t,n,function(r){return e[r]}.bind(null,n));return t},p.n=function(e){var r=e&&e.__esModule?function(){return e.default}:function(){return e};return p.d(r,"a",r),r},p.o=function(e,r){return Object.prototype.hasOwnProperty.call(e,r)},p.p="/";var l=this["webpackJsonpmy-app"]=this["webpackJsonpmy-app"]||[],a=l.push.bind(l);l.push=r,l=l.slice();for(var f=0;f<l.length;f++)r(l[f]);var i=a;t()}([])</script><script src="/static/js/2.a89711f7.chunk.js"></script><script src="/static/js/main.e3d6584e.chunk.js"></script></body></html>`

func BenchmarkZstdCompressor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		since := time.Now()
		byts := zstdc.Compress([]byte(html))
		b.Logf("orig size: %d, compressed size: %d", len(html), len(byts))
		b.Logf("compress took: %s", time.Now().Sub(since))
	}
}
