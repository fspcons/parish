# Firestore indexes

The API uses **Firestore native mode** (`cloud.google.com/go/firestore`). Composite indexes are required for queries that combine filters and `orderBy` (e.g. materials by `type` + `createdAt`).

Deploy indexes to your GCP project:

```bash
firebase deploy --only firestore:indexes --project YOUR_PROJECT_ID
```

Or copy `firestore.indexes.json` into your Firebase project’s `firestore.indexes.json` and deploy from there.

**Note:** Local development uses the Firestore emulator (`FIRESTORE_EMULATOR_HOST`); the emulator does not enforce missing indexes the same way production does.
