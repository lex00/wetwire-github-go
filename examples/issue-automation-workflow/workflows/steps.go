package workflows

import (
	"github.com/lex00/wetwire-github-go/actions/github_script"
)

// AutoLabelSteps labels issues based on content keywords.
var AutoLabelSteps = []any{
	github_script.GithubScript{
		Script: `
const issue = context.payload.issue;
const title = issue.title.toLowerCase();
const body = (issue.body || '').toLowerCase();
const content = title + ' ' + body;

const labelMap = {
  'bug': ['bug', 'error', 'fix', 'broken', 'crash'],
  'enhancement': ['feature', 'enhance', 'improve', 'add'],
  'documentation': ['docs', 'documentation', 'readme', 'typo'],
  'question': ['question', 'how to', 'help', 'confused']
};

const labelsToAdd = [];
for (const [label, keywords] of Object.entries(labelMap)) {
  if (keywords.some(keyword => content.includes(keyword))) {
    labelsToAdd.push(label);
  }
}

if (labelsToAdd.length > 0) {
  await github.rest.issues.addLabels({
    owner: context.repo.owner,
    repo: context.repo.repo,
    issue_number: issue.number,
    labels: labelsToAdd
  });
  console.log('Added labels:', labelsToAdd.join(', '));
} else {
  console.log('No matching labels found for this issue');
}
`,
	},
}

// RespondToCommentSteps responds to commands in issue comments.
var RespondToCommentSteps = []any{
	github_script.GithubScript{
		Script: `
const comment = context.payload.comment;
const body = comment.body.trim();

// Only respond to commands starting with /
if (!body.startsWith('/')) {
  console.log('Not a command, skipping');
  return;
}

const command = body.split(/\s+/)[0].toLowerCase();
const issueNumber = context.payload.issue.number;

const responses = {
  '/help': 'Available commands:\n- /help - Show this help message\n- /assign - Assign this issue to yourself\n- /close - Close this issue',
  '/assign': null, // Special handling below
  '/close': null   // Special handling below
};

if (command === '/assign') {
  await github.rest.issues.addAssignees({
    owner: context.repo.owner,
    repo: context.repo.repo,
    issue_number: issueNumber,
    assignees: [comment.user.login]
  });
  await github.rest.issues.createComment({
    owner: context.repo.owner,
    repo: context.repo.repo,
    issue_number: issueNumber,
    body: 'Assigned to @' + comment.user.login
  });
} else if (command === '/close') {
  await github.rest.issues.update({
    owner: context.repo.owner,
    repo: context.repo.repo,
    issue_number: issueNumber,
    state: 'closed'
  });
  await github.rest.issues.createComment({
    owner: context.repo.owner,
    repo: context.repo.repo,
    issue_number: issueNumber,
    body: 'Closed by @' + comment.user.login
  });
} else if (command === '/help') {
  await github.rest.issues.createComment({
    owner: context.repo.owner,
    repo: context.repo.repo,
    issue_number: issueNumber,
    body: responses['/help']
  });
} else {
  console.log('Unknown command:', command);
}
`,
	},
}

// EnforceReviewPolicySteps enforces review policies on PRs.
var EnforceReviewPolicySteps = []any{
	github_script.GithubScript{
		Script: `
const review = context.payload.review;
const pr = context.payload.pull_request;

// Only act on approved reviews
if (review.state !== 'approved') {
  console.log('Review state is not approved, skipping');
  return;
}

// Check if the reviewer is not the PR author (self-review prevention)
if (review.user.login === pr.user.login) {
  await github.rest.issues.createComment({
    owner: context.repo.owner,
    repo: context.repo.repo,
    issue_number: pr.number,
    body: 'Self-approval is not allowed. Please request a review from another team member.'
  });

  // Dismiss the self-approval
  await github.rest.pulls.dismissReview({
    owner: context.repo.owner,
    repo: context.repo.repo,
    pull_number: pr.number,
    review_id: review.id,
    message: 'Self-approval dismissed automatically'
  });
  return;
}

// Get all reviews for this PR
const reviews = await github.rest.pulls.listReviews({
  owner: context.repo.owner,
  repo: context.repo.repo,
  pull_number: pr.number
});

// Count unique approvals (excluding author)
const approvers = new Set();
for (const r of reviews.data) {
  if (r.state === 'APPROVED' && r.user.login !== pr.user.login) {
    approvers.add(r.user.login);
  }
}

console.log('Approved by', approvers.size, 'reviewer(s):', [...approvers].join(', '));

// Add label based on approval count
const labelToAdd = approvers.size >= 2 ? 'ready-to-merge' : 'needs-review';
const labelToRemove = approvers.size >= 2 ? 'needs-review' : 'ready-to-merge';

try {
  await github.rest.issues.removeLabel({
    owner: context.repo.owner,
    repo: context.repo.repo,
    issue_number: pr.number,
    name: labelToRemove
  });
} catch (e) {
  // Label might not exist, ignore
}

await github.rest.issues.addLabels({
  owner: context.repo.owner,
  repo: context.repo.repo,
  issue_number: pr.number,
  labels: [labelToAdd]
});
`,
	},
}
