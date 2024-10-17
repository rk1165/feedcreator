### feedcreator

- This project aims to turn any website to RSS Feed which we can then monitor using RSS readers.
- This [link](https://www.xul.fr/en-xml-rss.html) explains what are RSS feeds pretty well.

### How to use the app?

- To create RSS Feed we mainly need two things: `title` and `link`. There's an optional third thing, `description`,
  which can be skipped. A website consists of html pages which have elements like `<li>`, `<a>`, `<article>`, `<div>`
  etc.
- `li`, `a`, `article` are called **tags**. The elements **may** also have a `class` attribute associated with them.
- `class` attributes are used to apply `css` to a bunch of elements together. They also uniquely identifies
  elements on the webpage.
- To create RSS feed we need to identify such **common** elements on a webpage which will have `title` and `link`.
  Mostly, these are items appearing in a list format. Those two sub elements can also have class to identify them
  uniquely. With these three things we can create the main component of our RSS feed `<item>`.
- Consider the below HackerNews front page
  ![hackernews](docs/img/hn.png)
- This has 4 items in list format. In a day these items get updated, and we can use RSS feeds to track them.
- A single item is something like :
  ![item](docs/img/item.png)
- To identify the `element` associated with the list item we can right-click on the list title and select `inspect`
- The result is shown below:
  ![extractors](docs/img/extractors.png)
- The item here would be `span` with `class` attribute `titleline`
- The title and link element both will be `a` without any class attributes. The title will be the text content in the
  link.
- Once we have identified these items we need to fill the following form:

  ![form](docs/img/form.png)
    - *Title* is the title of the feed with which you would like to track the feed. Here it could be like
      `HackerNews Feed`
    - *URL* is the link of the feed. Here `https://news.ycombinator.com/newest`
    - *Description* is optional but could be something to describe the feed.
    - *Name* is something to uniquely identify the feed one is tracking. It could be a short name like `hn`
- Now we need to fill the extractor parameters. From the above example:
    - For item selector we have to use `span.titleline`
    - For title selector we have to use just `a`.
    - For link selector we have to use just `a`.
- After filling the form it should look like this:

  ![filled](docs/img/filled.png)

- Once we submit it we will get a page which will have the feed url with other details and will look like this:
  ![output](docs/img/output.png)
- Copy the `hn.xml` link - `http://localhost:8080/static/rss/hn.xml` and you can add it to any of the RSS Feed readers

### QuickStart

- `make init` to setup the database and download the dependencies.
- `make run` to start the application locally.
- `make build` to build the application locally.
- `make clean` to purge the contents of the database

### Design & Implementation

- This project is using [colly](https://github.com/gocolly/colly) library for scraping the webpage.
- The feeds are saved under `ui/static/rss` directory and there is a `feeds.db` table where we save feed metadata.
- Two functions are scheduled at a configurable interval in `main.go`. Using the metadata it rescans all **urls** to
  update or clean the feed items.
- For some websites the web page is not loaded completely before executing the JavaScript.
- For such pages we used `` to execute it and wait a second for it to load completely.
    - There is a fair chance that the page could still not be loaded completely in which case we won't be able to track
      it.

### Hosting

- There is a github workflow which creates a docker image for self hosting.

### Comments

- Currently, we rescan all the feeds at a fixed interval. We can optimize it to scan each site at a fixed interval.
- We are also not passing request headers like `If-Modified-Since` or checking response headers like `Last-Modified`.
- We can add a LLM button to find the extractors given a web page.
- For websites which only load completely after running JS we can create a separate slightly long-running process to
  load them.