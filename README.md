> Don't let an algorithm decide what you do or don't read. Receive a daily, weekly, or monthly email containing the latest articles from your favorite blogs.

## Workflow
sequenceDiagram
    User->>Web: Provide email address and <br>reading preferences
    Web->>Database: Save user preferences
    loop
        Note left of Cronjob: Run script via cronjob<br> every day/week/month
        Cronjob->>Database: Get user preferences
        Cronjob->>Blog(s): Get feeds
        Cronjob->>Cronjob: Create email
        Cronjob->>User: Send email containing recent publications
        Note right of User: Emails contains user <br>specific links
    end
    User->>Web: Edit preferences from email
    Web->>Database: Update subscriptions and/or user
    User->>Web: Unsubscribe
    Web->>Database: Delete user and associated subscriptions

## Deploy
[deploy.md](./deploy.md)