package github

// profileQuery pulls everything needed for the profile, stats and languages
// cards in one round trip. Repo pagination is handled by the caller if the
// user owns more than 100 repos.
const profileQuery = `
query($login: String!, $after: String) {
  user(login: $login) {
    id
    login
    name
    bio
    avatarUrl
    company
    location
    websiteUrl
    createdAt
    followers { totalCount }
    following { totalCount }
    pullRequests { totalCount }
    issues { totalCount }
    repositoriesContributedTo(
      first: 1
      contributionTypes: [COMMIT, PULL_REQUEST, ISSUE, PULL_REQUEST_REVIEW]
    ) { totalCount }
    contributionsCollection {
      totalCommitContributions
      totalIssueContributions
      totalPullRequestContributions
      totalPullRequestReviewContributions
      totalRepositoryContributions
      restrictedContributionsCount
      contributionCalendar { totalContributions }
    }
    repositories(
      first: 100
      after: $after
      ownerAffiliations: OWNER
      isFork: false
      orderBy: { field: STARGAZERS, direction: DESC }
    ) {
      totalCount
      pageInfo { hasNextPage endCursor }
      nodes {
        name
        stargazerCount
        forkCount
        primaryLanguage { name color }
        languages(first: 20, orderBy: { field: SIZE, direction: DESC }) {
          edges {
            size
            node { name color }
          }
        }
      }
    }
  }
}`

// commitHistoryQuery fetches commit timestamps in the default branch of one
// repo, filtered to commits authored by the target user. Used to build the
// productive-time heatmap.
const commitHistoryQuery = `
query($login: String!, $repo: String!, $userId: ID!, $since: GitTimestamp!, $after: String) {
  repository(owner: $login, name: $repo) {
    defaultBranchRef {
      target {
        ... on Commit {
          history(first: 100, after: $after, author: { id: $userId }, since: $since) {
            pageInfo { hasNextPage endCursor }
            nodes { committedDate }
          }
        }
      }
    }
  }
}`
