const API_ENDPOINT = 'https://api.rss2json.com/v1/api.json?rss_url=';
const STORAGE_KEY = 'rsssubs';

const RECOMMENDED = [
    { name: "A List Apart", blogUrl: "https://alistapart.com/", rssUrl: "https://alistapart.com/main/feed/" },
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
    { name: "Paul Graham", blogUrl: "https://paulgraham.com/articles.html", rssUrl: "http://www.aaronsw.com/2002/feeds/pgessays.rss" },
    { name: "Rachel by the Bay", blogUrl: "https://rachelbythebay.com/w/", rssUrl: "https://rachelbythebay.com/w/atom.xml" },
    { name: "Sam Altman's Blog", blogUrl: "https://blog.samaltman.com/", rssUrl: "https://blog.samaltman.com/posts.atom" },
    { name: "Seth's Blog", blogUrl: "https://seths.blog/", rssUrl: "https://feeds.feedblitz.com/sethsblog&x=1" },
    { name: "Staff Engineer", blogUrl: "https://staffeng.com/", rssUrl: "https://staffeng.com/rss" },
    { name: "Wait But Why", blogUrl: "https://waitbutwhy.com/", rssUrl: "https://waitbutwhy.com/feed" },
    { name: "web.dev", blogUrl: "https://web.dev/", rssUrl: "https://web.dev/feed.xml" }
];

// --- Storage helpers ---

function getSubs() {
    return JSON.parse(localStorage.getItem(STORAGE_KEY)) || [];
}

function saveSubs(subs) {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(subs));
}

// --- Articles ---

function isWithinLast7Days(dateStr) {
    return Date.now() - new Date(dateStr).getTime() <= 7 * 24 * 60 * 60 * 1000;
}

async function fetchFeedArticles(rssUrl) {
    try {
        const res = await fetch(`${API_ENDPOINT}${encodeURIComponent(rssUrl)}`);
        if (!res.ok) return [];
        const data = await res.json();
        if (!data.feed || !data.items?.length) return [];
        return data.items
            .filter(item => isWithinLast7Days(item.pubDate))
            .map(item => ({
                feedLink: data.feed.link,
                feedTitle: data.feed.title,
                author: item.author || '',
                date: new Date(item.pubDate),
                link: item.link,
                title: item.title
            }));
    } catch {
        return [];
    }
}

function createArticleEl(article) {
    const li = document.createElement('li');
    li.className = 'article';

    const titleEl = document.createElement('p');
    titleEl.className = 'article-title';
    const titleLink = document.createElement('a');
    titleLink.href = article.link;
    titleLink.target = '_blank';
    titleLink.rel = 'noopener noreferrer';
    titleLink.textContent = article.title;
    titleEl.appendChild(titleLink);

    const metaEl = document.createElement('div');
    metaEl.className = 'article-meta';

    const feedLink = document.createElement('a');
    feedLink.href = article.feedLink;
    feedLink.target = '_blank';
    feedLink.rel = 'noopener noreferrer';
    feedLink.textContent = article.feedTitle;
    metaEl.appendChild(feedLink);

    const addSep = () => {
        const sep = document.createElement('span');
        sep.className = 'meta-sep';
        sep.textContent = '·';
        metaEl.appendChild(sep);
    };

    if (article.author) {
        addSep();
        const authorEl = document.createElement('span');
        authorEl.textContent = article.author;
        metaEl.appendChild(authorEl);
    }

    addSep();
    const dateEl = document.createElement('span');
    dateEl.textContent = article.date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
    metaEl.appendChild(dateEl);

    li.appendChild(titleEl);
    li.appendChild(metaEl);
    return li;
}

async function loadArticles() {
    const list = document.getElementById('articleList');
    const countEl = document.getElementById('articleCount');
    list.innerHTML = '';

    const subs = getSubs();
    if (!subs.length) {
        const msg = document.createElement('li');
        msg.className = 'status-msg';
        msg.textContent = 'No feeds subscribed. Go to Feeds to add some.';
        list.appendChild(msg);
        countEl.textContent = '';
        return;
    }

    const loadingMsg = document.createElement('li');
    loadingMsg.className = 'status-msg';
    loadingMsg.textContent = 'Loading…';
    list.appendChild(loadingMsg);

    const results = await Promise.all(subs.map(feed => fetchFeedArticles(feed.rssUrl)));
    const articles = results.flat().sort((a, b) => b.date - a.date);

    list.innerHTML = '';

    if (!articles.length) {
        const msg = document.createElement('li');
        msg.className = 'status-msg';
        msg.textContent = 'No articles from the last 7 days.';
        list.appendChild(msg);
        countEl.textContent = '';
        return;
    }

    articles.forEach(a => list.appendChild(createArticleEl(a)));
    countEl.textContent = `(${articles.length})`;
}

// --- Subscriptions ---

function createFeedEl(feed, isSubscribed) {
    const li = document.createElement('li');
    li.className = 'feed-item';

    const info = document.createElement('div');
    info.className = 'feed-info';

    const nameEl = document.createElement('a');
    nameEl.className = 'feed-name';
    nameEl.href = feed.blogUrl || feed.rssUrl;
    nameEl.target = '_blank';
    nameEl.rel = 'noopener noreferrer';
    nameEl.textContent = feed.name || feed.rssUrl;

    const urlEl = document.createElement('span');
    urlEl.className = 'feed-url';
    urlEl.textContent = feed.rssUrl;

    info.appendChild(nameEl);
    info.appendChild(urlEl);

    const btn = document.createElement('button');
    btn.className = `btn-feed ${isSubscribed ? 'remove' : 'add'}`;
    btn.textContent = isSubscribed ? 'Remove' : 'Add';
    btn.dataset.rssUrl = feed.rssUrl;
    btn.dataset.action = isSubscribed ? 'unsubscribe' : 'subscribe';

    li.appendChild(info);
    li.appendChild(btn);
    return li;
}

function loadSubscriptions() {
    const subs = getSubs();
    const subList = document.getElementById('subscriptionList');
    const recList = document.getElementById('recommendedList');

    subList.innerHTML = '';
    recList.innerHTML = '';

    const subUrls = new Set(subs.map(f => f.rssUrl));

    document.getElementById('noSubs').classList.toggle('hidden', subs.length > 0);
    subs.forEach(feed => subList.appendChild(createFeedEl(feed, true)));

    RECOMMENDED.filter(f => !subUrls.has(f.rssUrl))
        .forEach(feed => recList.appendChild(createFeedEl(feed, false)));

    const countEl = document.getElementById('subCount');
    countEl.textContent = subs.length ? `(${subs.length})` : '';
}

function handleFeedAction(e) {
    const btn = e.target.closest('[data-action]');
    if (!btn) return;

    const { action, rssUrl } = btn.dataset;
    let subs = getSubs();

    if (action === 'subscribe') {
        const feed = RECOMMENDED.find(f => f.rssUrl === rssUrl);
        if (feed && !subs.some(f => f.rssUrl === rssUrl)) {
            saveSubs([...subs, feed]);
        }
    } else if (action === 'unsubscribe') {
        saveSubs(subs.filter(f => f.rssUrl !== rssUrl));
    }

    loadSubscriptions();
}

function addCustomFeed() {
    const input = document.getElementById('rssUrl');
    const errorEl = document.getElementById('addFeedError');
    const url = input.value.trim();

    errorEl.textContent = '';
    errorEl.classList.add('hidden');

    if (!url) return;

    const subs = getSubs();
    if (subs.some(f => f.rssUrl === url)) {
        errorEl.textContent = 'Already subscribed to this feed.';
        errorEl.classList.remove('hidden');
        return;
    }

    saveSubs([...subs, { name: url, blogUrl: url, rssUrl: url }]);
    input.value = '';
    loadSubscriptions();
}

// --- Navigation ---

function showSection(section) {
    const isArticles = section === 'articles';
    document.getElementById('articlesView').classList.toggle('hidden', !isArticles);
    document.getElementById('subscriptionsView').classList.toggle('hidden', isArticles);
    document.getElementById('showArticlesBtn').classList.toggle('active', isArticles);
    document.getElementById('showSubsBtn').classList.toggle('active', !isArticles);

    if (isArticles) loadArticles();
    else loadSubscriptions();
}

// --- Event listeners ---

document.getElementById('showArticlesBtn').addEventListener('click', () => showSection('articles'));
document.getElementById('showSubsBtn').addEventListener('click', () => showSection('subscriptions'));
document.getElementById('addFeedBtn').addEventListener('click', addCustomFeed);
document.getElementById('rssUrl').addEventListener('keydown', e => { if (e.key === 'Enter') addCustomFeed(); });
document.getElementById('rssUrl').addEventListener('input', () => {
    document.getElementById('addFeedError').classList.add('hidden');
});
document.getElementById('subscriptionList').addEventListener('click', handleFeedAction);
document.getElementById('recommendedList').addEventListener('click', handleFeedAction);

document.addEventListener('DOMContentLoaded', () => {
    loadSubscriptions();
    loadArticles();
});
