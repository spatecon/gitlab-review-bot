# GitLab Review Bot

A server application that observes GitLab merge requests (MR) and rotates reviews using custom Review&Approve policies.

<img alt="Go" src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white"> <img alt="GitHub" src="https://img.shields.io/github/license/spatecon/gitlab-review-bot?color=blue&style=for-the-badge">

## Situation

Have you ever been tagged for an ASAP review during the day? Switching the context of the developer wastes his energy
and reduces the entire team's productivity.

Down with `@here` tagging of whole team for MR review. Let the robot follow Review&Approve policy.

## Approach

What if we delegate supervising work to a machine?

- Reviewers assign on every MR: Why can't the bot randomly pick 2 of 5 developers of a team?
- Remind every team member about one's distinctive MR list. Right after daily meetings in Slack.
- Follow review rules: 2 approves move task into next status
- Gathering statistics: measure team performance over time, search for bottlenecks

## Solution

Review bot (written in Go)

- Easy to start: just list your teammates and projects
- Flexible: build your own review policy or choose between existing
- Lightweight: could be run on the cloud or a tiny dedicated server
- Open for extension: while the bot is written in go, it becomes easier to write custom Review&Approve policies using
  ready components and examples.

## Examples of Review&Approve policy

|          Policy          |                Reviews Selection                | Code Approval |
|:------------------------:|:-----------------------------------------------:|:-------------:|
| TeamLead is always right | random pick 2 devs and lead 👩‍💻🧑‍💻 + 🧙‍♂️️ |    1 lead     |
|     Developers riot      |          random pick 2 devs 👩‍💻🧑‍💻          |     1 dev     |
|  Reinventing Democracy   |          random pick 2 devs 👩‍💻👨‍💻          |    2 devs     |

---

## Roadmap

- [x] Reviewers rotation
- [x] Action on approved condition
- [ ] Slack reviews reminder (TODO)
    - [ ] Scheduled (e.g., after daily)
    - [ ] On user request
- [ ] Statistics gathering (TODO)
- [ ] Day off and vacation accounting (TODO)
- [ ] Jira task status integration (TODO)
- [x] Custom Review&Approve policies