package crawler

import (
	"fmt"
	"log"
	"regexp"
	"sync"
)

type crawler struct {
	EntryUrl       string
	MaxDepth       int
	CurDepth       int
	Visited        sync.Map
	NeedProcessing chan []string
}

func New(entryUrl string) crawler {
	c := crawler{
		EntryUrl:       entryUrl,
		MaxDepth:       2,
		CurDepth:       1,
		Visited:        sync.Map{},
		NeedProcessing: make(chan []string),
	}
	return c
}

// Crawl until there is something in c.NeedProcessing and all processPage runs have finished
func (c crawler) Crawl() {
	var wg sync.WaitGroup

	// start
	var n = 1 // pending URLs to process
	go func() { c.NeedProcessing <- []string{c.EntryUrl} }()

	for ; n > 0; n-- {
		list := <-c.NeedProcessing
		for _, url := range list {
			c.Visited.Store(url, true)
			wg.Add(1)
			log.Println("URL: %v %d", url, n)
			n++
			go c.processPage(url, &wg)
		}
	}

	// closer
	go func() {
		wg.Wait()
		//close(c.NeedProcessing)
	}()
}

func (c crawler) processPage(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	//// outside: check if hasn't been visited
	//client := http.Client{}
	//resp, err := client.Get(url)
	//if err != nil {
	//	log.Println(errors.Wrap(err, "processPage - can't create http client"))
	//}
	//
	//defer resp.Body.Close()
	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	log.Println(errors.Wrap(err, "processPage - can't read response"))
	//}

	var body = `
<!DOCTYPE html>
<html lang="en">
<meta charset="utf-8">
<meta name="description" content="Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#00ADD8">

  <title>The Go Project - The Go Programming Language</title>

<link href="https://fonts.googleapis.com/css?family=Work+Sans:600|Roboto:400,700" rel="stylesheet">
<link href="https://fonts.googleapis.com/css?family=Product+Sans&text=Supported%20by%20Google&display=swap" rel="stylesheet">
<link type="text/css" rel="stylesheet" href="/lib/godoc/style.css">

<link rel="search" type="application/opensearchdescription+xml" title="godoc" href="/opensearch.xml" />

<script>window.initFuncs = [];</script>

<script>
var _gaq = _gaq || [];
_gaq.push(["_setAccount", "UA-11222381-2"]);
window.trackPageview = function() {
  _gaq.push(["_trackPageview", location.pathname+location.hash]);
};
window.trackPageview();
window.trackEvent = function(category, action, opt_label, opt_value, opt_noninteraction) {
  _gaq.push(["_trackEvent", category, action, opt_label, opt_value, opt_noninteraction]);
};
</script>

<script src="/lib/godoc/jquery.js" defer></script>


<script src="/lib/godoc/playground.js" defer></script>

<script>var goVersion = "go1.16";</script>
<script src="/lib/godoc/godocs.js" defer></script>

<body class="Site">
<header class="Header js-header">
  <div class="Header-banner">
    Black Lives Matter.
    <a href="https://support.eji.org/give/153413/#!/donation/checkout"
       target="_blank"
       rel="noopener">Support the Equal Justice Initiative.</a>
  </div>
  <nav class="Header-nav Header-nav--wide">
    <a href="/"><img class="Header-logo" src="/lib/godoc/images/go-logo-blue.svg" alt="Go"></a>
    <button class="Header-menuButton js-headerMenuButton" aria-label="Main menu" aria-expanded="false">
      <div class="Header-menuButtonInner"></div>
    </button>
    <ul class="Header-menu">
      <li class="Header-menuItem"><a href="/doc/">Documents</a></li>
      <li class="Header-menuItem"><a href="/pkg/">Packages</a></li>
      <li class="Header-menuItem"><a href="/project/">The Project</a></li>
      <li class="Header-menuItem"><a href="/help/">Help</a></li>
      
        <li class="Header-menuItem"><a href="/blog/">Blog</a></li>
        
          <li class="Header-menuItem"><a href="https://play.golang.org/">Play</a></li>
        
      
    </ul>
  </nav>
</header>

<main id="page" class="Site-content wide">
<div class="container">


  <h1>
    The Go Project
    <span class="text-muted"></span>
  </h1>







<div id="nav"></div>




<img class="gopher" src="/doc/gopher/project.png" alt="" />

<div id="manual-nav"></div>

<p>
Go is an open source project developed by a team at
<a href="//google.com/">Google</a> and many
<a href="/CONTRIBUTORS">contributors</a> from the open source community.
</p>

<p>
Go is distributed under a <a href="/LICENSE">BSD-style license</a>.
</p>

<h3 id="announce"><a href="//groups.google.com/group/golang-announce">Announcements Mailing List</a></h3>
<p>
A low traffic mailing list for important announcements, such as new releases.
</p>
<p>
We encourage all Go users to subscribe to
<a href="//groups.google.com/group/golang-announce">golang-announce</a>.
</p>


<h2 id="go1">Version history</h2>

<h3 id="release"><a href="/doc/devel/release.html">Release History</a></h3>

<p>A <a href="/doc/devel/release.html">summary</a> of the changes between Go releases. Notes for the major releases:</p>

<ul>
	<li><a href="/doc/go1.16">Go 1.16</a> <small>(February 2021)</small></li>
	<li><a href="/doc/go1.15">Go 1.15</a> <small>(August 2020)</small></li>
	<li><a href="/doc/go1.14">Go 1.14</a> <small>(February 2020)</small></li>
	<li><a href="/doc/go1.13">Go 1.13</a> <small>(September 2019)</small></li>
	<li><a href="/doc/go1.12">Go 1.12</a> <small>(February 2019)</small></li>
	<li><a href="/doc/go1.11">Go 1.11</a> <small>(August 2018)</small></li>
	<li><a href="/doc/go1.10">Go 1.10</a> <small>(February 2018)</small></li>
	<li><a href="/doc/go1.9">Go 1.9</a> <small>(August 2017)</small></li>
	<li><a href="/doc/go1.8">Go 1.8</a> <small>(February 2017)</small></li>
	<li><a href="/doc/go1.7">Go 1.7</a> <small>(August 2016)</small></li>
	<li><a href="/doc/go1.6">Go 1.6</a> <small>(February 2016)</small></li>
	<li><a href="/doc/go1.5">Go 1.5</a> <small>(August 2015)</small></li>
	<li><a href="/doc/go1.4">Go 1.4</a> <small>(December 2014)</small></li>
	<li><a href="/doc/go1.3">Go 1.3</a> <small>(June 2014)</small></li>
	<li><a href="/doc/go1.2">Go 1.2</a> <small>(December 2013)</small></li>
	<li><a href="/doc/go1.1">Go 1.1</a> <small>(May 2013)</small></li>
	<li><a href="/doc/go1">Go 1</a> <small>(March 2012)</small></li>
</ul>

<h3 id="go1compat"><a href="/doc/go1compat">Go 1 and the Future of Go Programs</a></h3>
<p>
What Go 1 defines and the backwards-compatibility guarantees one can expect as
Go 1 matures.
</p>


<h2 id="resources">Developer Resources</h2>

<h3 id="source"><a href="https://golang.org/change">Source Code</a></h3>
<p>Check out the Go source code.</p>

<h3 id="discuss"><a href="//groups.google.com/group/golang-nuts">Discussion Mailing List</a></h3>
<p>
A mailing list for general discussion of Go programming.
</p>
<p>
Questions about using Go or announcements relevant to other Go users should be sent to
<a href="//groups.google.com/group/golang-nuts">golang-nuts</a>.
</p>

<h3 id="golang-dev"><a href="https://groups.google.com/group/golang-dev">Developer</a> and
<a href="https://groups.google.com/group/golang-codereviews">Code Review Mailing List</a></h3>
<p>The <a href="https://groups.google.com/group/golang-dev">golang-dev</a>
mailing list is for discussing code changes to the Go project.
The <a href="https://groups.google.com/group/golang-codereviews">golang-codereviews</a>
mailing list is for actual reviewing of the code changes (CLs).</p>

<h3 id="golang-checkins"><a href="https://groups.google.com/group/golang-checkins">Checkins Mailing List</a></h3>
<p>A mailing list that receives a message summarizing each checkin to the Go repository.</p>

<h3 id="build_status"><a href="//build.golang.org/">Build Status</a></h3>
<p>View the status of Go builds across the supported operating
systems and architectures.</p>


<h2 id="howto">How you can help</h2>

<h3><a href="//golang.org/issue">Reporting issues</a></h3>

<p>
If you spot bugs, mistakes, or inconsistencies in the Go project's code or
documentation, please let us know by
<a href="//golang.org/issue/new">filing a ticket</a>
on our <a href="//golang.org/issue">issue tracker</a>.
(Of course, you should check it's not an existing issue before creating
a new one.)
</p>

<p>
We pride ourselves on being meticulous; no issue is too small.
</p>

<p>
Security-related issues should be reported to
<a href="mailto:security@golang.org">security@golang.org</a>.<br>
See the <a href="/security">security policy</a> for more details.
</p>

<p>
Community-related issues should be reported to
<a href="mailto:conduct@golang.org">conduct@golang.org</a>.<br>
See the <a href="/conduct">Code of Conduct</a> for more details.
</p>

<h3><a href="/doc/contribute.html">Contributing code &amp; documentation</a></h3>

<p>
Go is an open source project and we welcome contributions from the community.
</p>
<p>
To get started, read these <a href="/doc/contribute.html">contribution
guidelines</a> for information on design, testing, and our code review process.
</p>
<p>
Check <a href="//golang.org/issue">the tracker</a> for
open issues that interest you. Those labeled
<a href="https://github.com/golang/go/issues?q=is%3Aopen+is%3Aissue+label%3A%22help+wanted%22">help wanted</a>
are particularly in need of outside help.
</p>


</div><!-- .container -->
</main><!-- #page -->
<footer>
  <div class="Footer Footer--wide">
    <img class="Footer-gopher" src="/lib/godoc/images/footer-gopher.jpg" alt="The Go Gopher">
    <ul class="Footer-links">
      <li class="Footer-link"><a href="/doc/copyright.html">Copyright</a></li>
      <li class="Footer-link"><a href="/doc/tos.html">Terms of Service</a></li>
      <li class="Footer-link"><a href="http://www.google.com/intl/en/policies/privacy/">Privacy Policy</a></li>
      <li class="Footer-link"><a href="http://golang.org/issues/new?title=x/website:" target="_blank" rel="noopener">Report a website issue</a></li>
    </ul>
    <a class="Footer-supportedBy" href="https://google.com">Supported by Google</a>
  </div>
</footer>

<script>
(function() {
  var ga = document.createElement("script"); ga.type = "text/javascript"; ga.async = true;
  ga.src = ("https:" == document.location.protocol ? "https://ssl" : "http://www") + ".google-analytics.com/ga.js";
  var s = document.getElementsByTagName("script")[0]; s.parentNode.insertBefore(ga, s);
})();
</script>


`

	re := regexp.MustCompile("href=\"(https:\\/\\/golang.org)*(?P<relativePath>\\/\\w+)\\/*\"")
	matches := re.FindAllSubmatch([]byte(body), 15)
	var fetchedUrls []string
	for _, match := range matches {
		for j, name := range re.SubexpNames() {
			if name == "relativePath" {
				newUrl := string(match[j])
				// If the url is relative, prepend it with the EntryUrl
				if len(newUrl) > 0 && newUrl[0] == '/' {
					newUrl = fmt.Sprintf("%v%v", c.EntryUrl, newUrl)
				}
				if _, ok := c.Visited.Load(newUrl); !ok {
					// queue for processing
					log.Println(fmt.Sprintf("queue for processing: %v", newUrl))
					fetchedUrls = append(fetchedUrls, newUrl)
				}
			}
		}
	}

	c.NeedProcessing <- fetchedUrls
	c.CurDepth++
	//if c.CurDepth >= c.MaxDepth {
	//	close(c.NeedProcessing)
	//}
	log.Println(fmt.Sprintf("processed: %v", url))

	// save to disk
	//fullPath, err := os.Getwd()
	//if err != nil {
	//	log.Println(errors.Wrap(err, "can not get full path"))
	//}
	//log.Println(fullPath)
	//dir := "root"
	//os.Mkdir(dir, 0777)
	//fileName := path.Join(dir, url)
	//if err := ioutil.WriteFile(fileName, []byte(body), 0666); err != nil {
	//	log.Println(errors.Wrap(err, "can not save file to disk"))
	//}
}
