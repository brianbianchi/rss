const apiEndpoint = 'https://api.rss2json.com/v1/api.json?rss_url=';
const localStorageKey = "bbfeedsubs";
const recommended = [
    { name: "A List Apart", blogUrl: "https://alistapart.com/", rssUrl: "https://alistapart.com/main/feed/" },
    { name: "Paul Graham", blogUrl: "https://paulgraham.com/articles.html", rssUrl: "http://www.aaronsw.com/2002/feeds/pgessays.rss" },
    { name: "Bakadesuyo", blogUrl: "https://bakadesuyo.com/", rssUrl: "https://bakadesuyo.com/feed/" },
    { name: "Coding Horror", blogUrl: "https://blog.codinghorror.com/", rssUrl: "https://blog.codinghorror.com/rss/" },
    { name: "CSS-Tricks", blogUrl: "https://css-tricks.com/", rssUrl: "https://css-tricks.com/feed/" },
    { name: "DEV Community", blogUrl: "https://dev.to/", rssUrl: "https://dev.to/feed/" },
    { name: "Farnam Street", blogUrl: "https://fs.blog/blog/", rssUrl: "https://fs.blog/blog/feed/" },
    { name: "Hacker News", blogUrl: "https://news.ycombinator.com/", rssUrl: "https://news.ycombinator.com/rss" },
    { name: "Joel on Software", blogUrl: "https://www.joelonsoftware.com/", rssUrl: "https://www.joelonsoftware.com/feed/" },
    { name: "Lethain's Blog", blogUrl: "https://lethain.com/", rssUrl: "https://lethain.com/feeds.xml" },
    { name: "Levels.io Blog", blogUrl: "https://levels.io/", rssUrl: "https://levels.io/rss/" },
    { name: "LogRocket Blog", blogUrl: "https://blog.logrocket.com/", rssUrl: "https://blog.logrocket.com/feed/" },
    { name: "Melting Asphalt", blogUrl: "https://www.meltingasphalt.com/", rssUrl: "https://feeds.feedburner.com/MeltingAsphalt" },
    { name: "Netflix TechBlog", blogUrl: "https://netflixtechblog.com/", rssUrl: "https://netflixtechblog.com/feed" },
    { name: "OpenAI Blog", blogUrl: "https://openai.com/blog", rssUrl: "https://openai.com/blog/rss/" },
    { name: "Rachel by the Bay", blogUrl: "https://rachelbythebay.com/w/", rssUrl: "https://rachelbythebay.com/w/atom.xml" },
    { name: "Reddit (All)", blogUrl: "https://www.reddit.com/", rssUrl: "https://www.reddit.com/.rss" },
    { name: "Reddit (Programming)", blogUrl: "https://www.reddit.com/r/programming/", rssUrl: "https://www.reddit.com/r/programming/.rss" },
    { name: "Reddit (ProgrammerHumor)", blogUrl: "https://www.reddit.com/r/ProgrammerHumor/", rssUrl: "https://www.reddit.com/r/ProgrammerHumor/.rss" },
    { name: "Sam Altman's Blog", blogUrl: "https://blog.samaltman.com/", rssUrl: "https://blog.samaltman.com/posts.atom" },
    { name: "Seth's Blog", blogUrl: "https://seths.blog/", rssUrl: "https://feeds.feedblitz.com/sethsblog&x=1" },
    { name: "Staff Engineer", blogUrl: "https://staffeng.com/", rssUrl: "https://staffeng.com/rss" },
    { name: "Wait But Why", blogUrl: "https://waitbutwhy.com/", rssUrl: "https://waitbutwhy.com/feed" },
    { name: "web.dev", blogUrl: "https://web.dev/", rssUrl: "https://web.dev/feed.xml" }
];

let articles = [];

function isWithinLast7Days(date) {
    const currentDate = new Date();
    const sevenDaysAgo = new Date().setDate(currentDate.getDate() - 7);
    const inputDate = new Date(date);
    return inputDate >= sevenDaysAgo && inputDate <= currentDate;
}

function createArticleItem(article) {
    const articleItem = document.createElement('li');
    articleItem.className = 'article';
    articleItem.innerHTML = `
        <a href="${article.link}" target="_blank" rel="noopener noreferrer">${article.title}</a>
        <br />
        <a href="${article.feedLink}" target="_blank" rel="noopener noreferrer">${article.feedTitle}</a>
        <br />
        ${article.author} ${new Date(article.date).toLocaleDateString("en-US")} <hr />
    `;
    return articleItem;
}

async function fetchArticles(rssUrl) {
    try {
        const response = await fetch(`${apiEndpoint}${rssUrl}`);
        const data = await response.json();
        if (response.status === 422) {
            console.error(`Error fetching articles for ${rssUrl}:`, data.message);
            return;
        }
        if (!data.feed || !data.items || data.items.length === 0) {
            console.warn(`No feed data or items found for ${rssUrl}. Response:`, data);
            return;
        }
        const feedUrl = data.feed.url;
        const feedLink = data.feed.link;
        const feedTitle = data.feed.title;
        data.items.forEach(article => {
            console.log(article.pubDate);

            if (!isWithinLast7Days(article.pubDate)) return;
            articles.push({
                feedLink,
                feedUrl,
                feedTitle,
                author: article.author,
                date: new Date(article.pubDate),
                link: article.link,
                title: article.title
            })
        });
    } catch (error) {
        console.error(error);
    }
}

function displayArticles() {
    const articlesDiv = document.getElementById('articles');
    articlesDiv.innerHTML = '';
    if (articles.length) {
        articlesDiv.innerHTML += `<h2>Articles</h2>
                                    <ul id="articleList"></ul>`;
        articles.forEach(article => {
            document.getElementById('articleList')
                .appendChild(createArticleItem(article));
        });
    } else {
        articlesDiv.innerHTML = `<p>No articles found. Add some subscriptions!</p>`;
    }
    document.getElementById('articleCount').textContent = `(${articles.length})`;
}

function createRssItem(feedData, isSubscribed) {
    const rssItem = document.createElement('li');
    rssItem.className = 'rss-item';

    let rssUrl = typeof feedData === 'string' ? feedData : feedData.rssUrl;
    let displayName = typeof feedData === 'string' ? feedData : (feedData.name || feedData.rssUrl);
    let displayLink = typeof feedData === 'string' ? null : (feedData.blogUrl || feedData.rssUrl);

    let displayHtml = '';
    if (displayLink && displayLink !== rssUrl) {
        displayHtml = `<a href="${displayLink}" target="_blank" rel="noopener noreferrer">${displayName}</a>`;
    } else {
        displayHtml = `<span>${displayName}</span>`;
    }

    rssItem.innerHTML = `
        ${displayHtml}
        <br />
        <small>${rssUrl}</small>
        ${isSubscribed ?
            `<button onclick="unsubscribe('${rssUrl}', this)">Unsubscribe</button>` :
            `<button onclick="subscribe('${rssUrl}', this)">Subscribe</button>`}
    `;
    return rssItem;
}

function toggleSubWarning(subs) {
    const warning = document.getElementById('noSubs');
    warning.style.display = subs.length ? 'none' : 'block';
}

function toggleContent(header) {
    const content = header.nextElementSibling;
    content.style.display = content.style.display === 'block' ? 'none' : 'block';
}

async function addCustomRssUrl() {
    const urlInput = document.getElementById('rssUrl');
    const url = urlInput.value.trim();

    if (url) {
        let storedFeeds = JSON.parse(localStorage.getItem(localStorageKey)) || [];
        if (!storedFeeds.some(feed => feed.rssUrl === url)) {
            const newFeed = { name: "Custom Feed", blogUrl: url, rssUrl: url };
            storedFeeds.push(newFeed);
            localStorage.setItem(localStorageKey, JSON.stringify(storedFeeds));

            document.getElementById('subscriptionList')
                .appendChild(createRssItem(newFeed, true));
            urlInput.value = '';

            toggleSubWarning(storedFeeds);
            updateSubscriptionCount();
            await loadArticles();
        } else {
            alert('This RSS feed is already subscribed!');
        }
    }
}

async function subscribe(url, button) {
    let storedFeeds = JSON.parse(localStorage.getItem(localStorageKey)) || [];
    const subscribedFeed = recommended.find(feed => feed.rssUrl === url);

    if (subscribedFeed && !storedFeeds.some(feed => feed.rssUrl === subscribedFeed.rssUrl)) {
        storedFeeds.push(subscribedFeed);
        localStorage.setItem(localStorageKey, JSON.stringify(storedFeeds));

        document.getElementById('recommendedList')
            .removeChild(button.parentElement);

        document.getElementById('subscriptionList')
            .appendChild(createRssItem(subscribedFeed, true));
        toggleSubWarning(storedFeeds);
        updateSubscriptionCount();
        await loadArticles();
    }
}

function unsubscribe(url, button) {
    let storedFeeds = JSON.parse(localStorage.getItem(localStorageKey)) || [];
    const newStoredFeeds = storedFeeds.filter(feed => feed.rssUrl !== url);
    localStorage.setItem(localStorageKey, JSON.stringify(newStoredFeeds));

    document.getElementById('subscriptionList').removeChild(button.parentElement);
    toggleSubWarning(newStoredFeeds);
    updateSubscriptionCount();
    loadArticles();
    loadSubscriptions();
}

function showSection(sectionId) {
    document.getElementById('articlesView').classList.add('hidden');
    document.getElementById('subscriptionsView').classList.add('hidden');

    document.getElementById('showArticlesBtn').classList.remove('active');
    document.getElementById('showSubsBtn').classList.remove('active');

    document.getElementById(sectionId).classList.remove('hidden');
    if (sectionId === 'articlesView') {
        document.getElementById('showArticlesBtn').classList.add('active');
        loadArticles();
    } else {
        document.getElementById('showSubsBtn').classList.add('active');
        loadSubscriptions();
    }
}

async function loadArticles() {
    articles = [];
    const storedFeeds = JSON.parse(localStorage.getItem(localStorageKey)) || [];
    for (const feed of storedFeeds) {
        await fetchArticles(feed.rssUrl);
    }
    articles.sort((a, b) => new Date(a.date) - new Date(b.date));
    displayArticles();
}

function loadSubscriptions() {
    const subscriptionList = document.getElementById('subscriptionList');
    subscriptionList.innerHTML = '';
    const recommendedList = document.getElementById('recommendedList');
    recommendedList.innerHTML = '';

    let storedFeeds = JSON.parse(localStorage.getItem(localStorageKey)) || [];
    storedFeeds.forEach(feed => {
        subscriptionList.appendChild(createRssItem(feed, true));
    });

    const storedRssUrls = storedFeeds.map(feed => feed.rssUrl);

    recommended.forEach(recFeed => {
        if (!storedRssUrls.includes(recFeed.rssUrl)) {
            recommendedList.appendChild(createRssItem(recFeed, false));
        }
    });
    toggleSubWarning(storedFeeds);
    updateSubscriptionCount();
}

function updateSubscriptionCount() {
    const subs = JSON.parse(localStorage.getItem(localStorageKey)) || [];
    document.getElementById('subCount').textContent = `(${subs.length})`;
}

document.getElementById('showArticlesBtn').addEventListener('click', () => showSection('articlesView'));
document.getElementById('showSubsBtn').addEventListener('click', () => showSection('subscriptionsView'));

document.addEventListener('DOMContentLoaded', () => {
    showSection('articlesView');
    loadSubscriptions();
});