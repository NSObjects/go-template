<!--See the general documentation guidelines https://docs.gitlab.com/ee/development/documentation -->

<!-- Mention "documentation" or "docs" in the MR title -->

<!-- Use this description template for new docs or updates to existing docs. For changing documentation location use the "Change documentation location" template -->

## 这个MR是干嘛的?

<!-- Briefly describe what this MR is about -->

## 相关的 issues

<!-- Mention the issue(s) this MR closes or is related to -->

Closes 

## checklist

(完成这个MR需要哪些步骤)

- [ ] [Apply the correct labels and milestone](https://docs.gitlab.com/ee/development/documentation/workflow.html#2-developer-s-role-in-the-documentation-process)
- [ ] Crosslink the document from the higher-level index
- [ ] Crosslink the document from other subject-related docs
- [ ] Feature moving tiers? Make sure the change is also reflected in [`features.yml`](https://gitlab.com/gitlab-com/www-gitlab-com/blob/master/data/features.yml)
- [ ] Correctly apply the product [badges](https://docs.gitlab.com/ee/development/documentation/styleguide.html#product-badges) and [tiers](https://docs.gitlab.com/ee/development/documentation/styleguide.html#gitlab-versions-and-tiers)
- [ ] [Port the MR to EE (or backport from CE)](https://docs.gitlab.com/ee/development/documentation/index.html#cherry-picking-from-ce-to-ee): _always recommended, required when the `ee-compat-check` job fails_

## Review checklist

(提醒Review人员需要去检查那部分代码)

- [ ] Your team's review (required)
- [ ] PM's review (recommended, but not a blocker)
- [ ] Technical writer's review (required)
- [ ] Merge the EE-MR first, CE-MR afterwards
