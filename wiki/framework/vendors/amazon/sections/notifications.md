# Amazon SP-API: Notifications

Last verified: 2026-04-27

## What this section covers

- Notifications API setup and destination management
- SQS and EventBridge delivery models
- App notification template and feedback flows

## Representative docs

- [Notifications API](https://developer-docs.amazon.com/sp-api/docs/notifications-api)
- [Set up Notifications with SQS](https://developer-docs.amazon.com/sp-api/docs/set-up-notifications-with-amazon-sqs)
- [Set up Notifications with EventBridge](https://developer-docs.amazon.com/sp-api/docs/set-up-notifications-with-amazon-eventbridge)
- [Tutorial: Grant SP-API Permission to SQS Queue](https://developer-docs.amazon.com/sp-api/docs/tutorial-grant-permission-to-sqs-queue)
- [Notification Type Values](https://developer-docs.amazon.com/sp-api/docs/notification-type-values)
- [View Order Notifications](https://developer-docs.amazon.com/sp-api/docs/view-order-notifications)
- [Record App Notification Feedback](https://developer-docs.amazon.com/sp-api/docs/record-app-notification-feedback)
- [Sample Sandbox Notification Templates](https://developer-docs.amazon.com/sp-api/docs/sample-sandbox-notification-templates)

## MPC notes

- SQS standard delivery can duplicate/reorder events; enforce dedupe and idempotent handlers.
- Keep polling reconciliation for critical entities as a fallback path.

