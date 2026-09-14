CASE CLICKER - DOWNLOAD WEBSITE
================================

Everything in this folder IS the website. Upload the whole folder as-is and it
works - the download buttons point at the two .exe files sitting next to
index.html, so there is nothing to configure.

  index.html                 the page
  img/                       screenshots + icon
  favicon.ico                tab icon
  Case Clicker Setup.exe     the installer (what the main button serves)
  Case Clicker.exe           the portable build (the secondary link)

HOSTING
-------
Netlify (easiest)
  Go to app.netlify.com/drop and drag this whole folder onto the page. You get
  a live URL in about 10 seconds. Free, and it allows .exe files.

Neocities
  Heads up: FREE Neocities accounts cannot host .exe files - only html/css/js/
  images and similar. Two options:
    a) Upgrade to Supporter, which allows any file type, or
    b) Upload only index.html, img/ and favicon.ico to Neocities, and repoint
       the two download links at something that can host the .exe. A GitHub
       Release works well and is free.

GitHub Pages
  Works fine, and will serve the .exe too. Push this folder to a repo and
  enable Pages on it.

UPDATING THE VERSION LATER
--------------------------
Three spots in index.html mention the version or size:
  - the .meta line under the main download button
  - the footer line "Case Clicker 1.3.0 - made by Foxyyyy"
  - the "Free, 2.1 MB, no strings." line above the bottom button
