# Claude/GitHub Integration Guide

This document explains the Claude Code GitHub integration that has been set up for this project, providing AI-assisted development capabilities directly within the GitHub workflow.

## 🎯 Overview

The Claude/GitHub integration provides two main capabilities:

1. **On-demand assistance** via `@claude` mentions in issues and pull requests
2. **Automatic code reviews** for all pull requests

This integration allows Claude to act as an AI pair programmer, reviewer, and assistant directly within your GitHub workflow, helping with code quality, bug detection, and development guidance.

## 📁 Current Workflow Files

The project includes four GitHub Actions workflow files that work together to provide a complete development experience:

### 1. `claude.yml` - On-Demand Claude Assistance

**Location**: `.github/workflows/claude.yml`

**Purpose**: Provides interactive Claude assistance when mentioned with `@claude`

**Triggers**:
- Issue comments containing `@claude`
- Pull request review comments containing `@claude`
- Pull request reviews containing `@claude`
- Issues opened with `@claude` in title or body

**Capabilities**:
- Answer questions about the codebase
- Help debug issues
- Suggest code improvements
- Explain complex code patterns
- Assist with Go best practices
- Help with CI/CD troubleshooting

**Permissions**:
- `contents: read` - Read repository files
- `pull-requests: read` - Access PR information
- `issues: read` - Access issue information
- `actions: read` - Read CI results on PRs
- `id-token: write` - GitHub OIDC authentication

### 2. `claude-code-review.yml` - Automatic Code Reviews

**Location**: `.github/workflows/claude-code-review.yml`

**Purpose**: Automatically reviews every pull request with AI-powered feedback

**Triggers**:
- Pull requests opened
- Pull requests synchronized (new commits pushed)

**Review Focus Areas**:
- Code quality and best practices
- Potential bugs or issues
- Performance considerations
- Security concerns
- Test coverage

**Benefits**:
- Consistent review standards
- Catches common issues early
- Educational feedback for developers
- Reduces reviewer workload

### 3. `ci.yml` - Continuous Integration (Build & Test)

**Location**: `.github/workflows/ci.yml`

**Purpose**: Comprehensive build, test, and quality checks for the Go microservices

**Triggers**:
- Push to `main` or `develop` branches (when `/api` files change)
- Pull requests to `main` or `develop` branches (when `/api` files change)

**What It Does**:
- **Multi-module detection**: Automatically detects which Go modules changed
- **Go testing**: Runs unit tests with PostgreSQL database
- **Integration testing**: Database migration and integration tests
- **Security scanning**: Basic vulnerability checks with Trivy
- **Container builds**: Builds and pushes Docker images to GitHub Container Registry
- **Quality gates**: Ensures all checks pass before allowing merges

**Key Features**:
- **Path-based triggers**: Only runs when API code changes
- **PostgreSQL service**: Spins up test database automatically
- **Multi-module support**: Handles both `billing` and `billing-dbmigrations`
- **Caching**: Caches Go modules for faster builds
- **Artifact uploads**: Saves test coverage reports

**Jobs**:
1. `detect-changes` - Determines which modules changed
2. `test-billing` - Tests the main billing service
3. `test-migrations` - Tests database migration tools
4. `security-scan` - Basic security vulnerability scanning
5. `build-images` - Builds and pushes container images
6. `quality-gate` - Final check that all required jobs passed

### 4. `dependabot-ci.yml` - Simplified Dependency CI

**Location**: `.github/workflows/dependabot-ci.yml`

**Purpose**: Lightweight CI specifically for Dependabot dependency update PRs

**Triggers**:
- Pull requests (only when authored by `dependabot[bot]`)

**What It Does**:
- **Dependency verification**: Downloads and verifies new dependencies
- **Basic compilation**: Ensures code still compiles with updates
- **Vulnerability checks**: Scans for known security issues
- **Compatibility testing**: Basic tests without full database setup

**Why Separate**:
- Dependabot PRs have **limited permissions** and can't access full CI secrets
- **Faster execution** for simple dependency updates
- **Reduced complexity** - focuses only on compatibility checks
- **Better success rate** - avoids issues with database setup in Dependabot context

**Benefits**:
- Quick feedback on dependency safety
- Automatic vulnerability detection
- Simplified approval process for minor updates

## 🔄 Workflow Orchestration

### How the Workflows Work Together

The four workflows coordinate to provide a seamless development experience:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Developer     │    │   Pull Request   │    │  Main Branch    │
│   Creates PR    │───▶│   Triggered      │───▶│   Protected     │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  Parallel Jobs   │
                    │                  │
                    │ • ci.yml         │
                    │ • claude-review  │
                    │ • dependabot-ci  │
                    │   (if applicable)│
                    └──────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  All Checks Pass │
                    │                  │
                    │ ✅ Tests pass    │
                    │ ✅ Builds work   │
                    │ ✅ Security OK   │
                    │ ✅ Claude review │
                    └──────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  Merge Allowed   │
                    │                  │
                    │ Branch protected │
                    │ by settings.yml  │
                    └──────────────────┘
```

### Typical Development Flow

1. **Create Feature Branch**: `git checkout -b feature/new-endpoint`
2. **Make Changes**: Edit Go code in `/api/billing` or `/api/billing-dbmigrations`
3. **Push Changes**: `git push origin feature/new-endpoint`
4. **Create PR**: GitHub automatically triggers workflows
5. **Automated Checks Run**:
   - `ci.yml` builds, tests, and scans your code
   - `claude-code-review.yml` provides AI feedback
   - `dependabot-ci.yml` runs if it's a dependency update
6. **Review Process**:
   - Check CI results for any failures
   - Read Claude's review comments
   - Address any issues or ask `@claude` for help
7. **Merge**: Once all checks pass and review is approved
8. **Cleanup**: Branch is automatically deleted

### Branch Protection Integration

The workflows integrate with branch protection rules defined in `.github/settings.yml`:

#### Required Status Checks
```yaml
required_status_checks:
  contexts:
    - "Continuous Integration"           # ci.yml must pass
    - "claude-review"                   # Claude review must complete
    - "dependabot-check / dependabot-summary"  # Dependabot CI (if applicable)
```

#### Protection Rules
- **No direct pushes** to `main` branch
- **Pull requests required** for all changes
- **Status checks must pass** before merge
- **Admin override allowed** (for learning purposes)

## 🚀 How to Use the Integration

### On-Demand Assistance

You can request Claude's help by mentioning `@claude` in:

#### Issue Comments
```markdown
I'm having trouble with the database connection in the billing service.
@claude can you help me debug this error:

Error: pq: database "billing" does not exist
```

#### Pull Request Comments
```markdown
@claude please review this function and suggest improvements:

func (c *Client) Validate() error {
    // implementation here
}
```

#### Pull Request Reviews
```markdown
@claude this PR looks good overall, but can you check if there are any 
security issues I might have missed?
```

#### New Issues
Create an issue with `@claude` in the title or description:
```markdown
Title: @claude Help with Go testing patterns

I need guidance on implementing table-driven tests for the billing service...
```

### Automatic Code Reviews

Every pull request automatically receives a code review from Claude focusing on:

- **Go best practices**: Proper error handling, idiomatic code patterns
- **Security**: Input validation, SQL injection prevention, secret handling
- **Performance**: Efficient algorithms, memory usage, database queries
- **Testing**: Test coverage, edge cases, test quality
- **Documentation**: Code comments, API documentation

## ⚙️ Configuration Options

Both workflows support extensive customization through optional parameters:

### Model Selection
```yaml
# Default: Claude Sonnet 4
model: "claude-opus-4-20250514"  # Use Claude Opus 4 for more complex tasks
```

### Custom Trigger Phrases
```yaml
# Change from @claude to a custom phrase
trigger_phrase: "/claude"
# or
trigger_phrase: "/ai-review"
```

### Project-Specific Instructions
```yaml
custom_instructions: |
  Follow our Go coding standards:
  - Use explicit error handling
  - Write table-driven tests
  - Document all exported functions
  - Use structured logging with logrus
  - Follow our REST API patterns
```

### Allowed Tools
```yaml
# Allow Claude to run specific commands
allowed_tools: "Bash(go test ./...),Bash(go mod tidy),Bash(golangci-lint run)"
```

### Environment Variables
```yaml
claude_env: |
  GO_ENV: development
  DATABASE_URL: postgresql://localhost:5432/billing_test
```

### Conditional Reviews
```yaml
# Only review PRs from external contributors
if: |
  github.event.pull_request.author_association == 'FIRST_TIME_CONTRIBUTOR' ||
  github.event.pull_request.author_association == 'CONTRIBUTOR'
```

### File-Specific Reviews
```yaml
# Only review specific file types
paths:
  - "api/**/*.go"
  - "*.md"
  - ".github/workflows/*.yml"
```

## 🔐 Security and Permissions

### What Claude Can Access

**Read Access**:
- Repository source code
- Pull request diffs
- Issue content and comments
- CI/CD workflow results
- Public repository metadata

**Write Access**:
- Post comments on issues and PRs
- Update PR reviews
- Create issue comments

**Cannot Access**:
- Repository secrets
- Private personal information
- Other repositories (unless explicitly configured)
- Administrative functions

### Authentication

The integration uses OAuth tokens stored in GitHub Secrets:
- `CLAUDE_CODE_OAUTH_TOKEN` - Authenticates Claude with GitHub
- Tokens are managed securely by GitHub
- Permissions are scoped to specific repository actions

### Data Privacy

- Claude only accesses repository content when triggered
- No persistent storage of repository data
- All interactions logged in GitHub Actions
- Follows Anthropic's data handling policies

## 💡 Usage Examples

### 1. Code Review Request
**Issue Comment**:
```markdown
@claude can you review this database connection code and check for potential issues?

```go
func NewDatabase(config *Config) (*sql.DB, error) {
    db, err := sql.Open("postgres", config.DatabaseURL)
    if err != nil {
        return nil, err
    }
    return db, nil
}
```

**Expected Response**: Claude will analyze the code and suggest improvements like connection pooling, ping testing, and proper error handling.

### 2. Debugging Assistance
**PR Comment**:
```markdown
@claude I'm getting this test failure and I'm not sure why:

```
=== RUN   TestClientValidation
    client_test.go:25: expected error for invalid email, got nil
--- FAIL: TestClientValidation (0.00s)
```

Can you help me understand what's wrong with my test?
```

**Expected Response**: Claude will examine the test code, identify the issue, and suggest fixes.

### 3. Architecture Guidance
**Issue Description**:
```markdown
@claude I need to add user authentication to the billing service. 

Current architecture:
- REST API with Gorilla Mux
- PostgreSQL database
- Docker containers
- OpenShift deployment

What's the best approach for adding JWT-based authentication?
```

**Expected Response**: Claude will suggest implementation patterns, security considerations, and code examples.

### 4. Performance Optimization
**PR Review**:
```markdown
@claude this PR adds a new endpoint that queries the database. Can you review it for performance issues?

The endpoint will be called frequently and needs to be fast.
```

**Expected Response**: Claude will analyze the database queries, suggest optimizations, and recommend caching strategies.

## 🔧 Integration with Existing CI/CD

The Claude integration complements the existing CI/CD workflows:

### Workflow Coordination
- **CI workflows** run automated tests and builds
- **Claude reviews** provide human-like code analysis
- **Security scans** detect vulnerabilities
- **Dependabot** manages dependency updates

### Typical PR Flow
1. **Developer** creates pull request
2. **CI workflows** run tests, linting, security scans
3. **Claude** automatically reviews code quality
4. **Human reviewers** focus on business logic and architecture
5. **Integration tests** run in staging environment
6. **Deployment** to production after approval

### Enhanced Debugging
When CI fails:
1. Check CI logs for errors
2. Mention `@claude` with error details
3. Claude analyzes logs and suggests fixes
4. Apply fixes and re-run CI

## 🚨 Troubleshooting

### Workflow-Specific Issues

#### 1. CI Workflow (`ci.yml`) Problems
**Symptoms**: Tests fail, builds don't complete, security scans error

**Common Causes**:
- PostgreSQL service not ready (wait longer)
- Go module dependencies outdated (`go mod tidy`)
- Missing environment variables
- Path filters not matching changed files

**Solutions**:
```bash
# Check workflow triggers
git log --oneline -5  # See what changes triggered workflow

# Update dependencies
cd api/billing && go mod tidy
cd api/billing-dbmigrations && go mod tidy

# Test locally
cd api/billing && go test ./...
```

#### 2. Dependabot CI Issues
**Symptoms**: Dependabot PRs fail, dependency conflicts

**Common Causes**:
- Dependency incompatibilities
- Breaking changes in new versions
- Go module conflicts

**Solutions**:
- Check Dependabot PR description for breaking changes
- Test locally: `go get [dependency]@[version]`
- Use `@claude` to analyze the specific dependency update

#### 3. Claude Workflows Not Triggering
**Symptoms**: No automatic reviews, `@claude` mentions ignored

**Common Causes**:
- OAuth token expired or missing
- Workflow file syntax errors
- Trigger conditions not met

**Solutions**:
```bash
# Check workflow syntax
gh workflow list
gh workflow view claude.yml

# Verify OAuth token exists
gh secret list | grep CLAUDE

# Check recent workflow runs
gh run list --workflow=claude.yml
```

### Common Issues

#### 1. Claude Not Responding
**Problem**: Mentioned `@claude` but no response

**Solutions**:
- Check if `CLAUDE_CODE_OAUTH_TOKEN` secret exists
- Verify workflow file syntax
- Check GitHub Actions logs
- Ensure correct repository permissions

#### 2. Incomplete Reviews
**Problem**: Claude's review seems superficial

**Solutions**:
- Add more specific `custom_instructions`
- Use `direct_prompt` for targeted feedback
- Enable specific tools with `allowed_tools`
- Consider upgrading to Claude Opus 4

#### 3. Too Many Reviews
**Problem**: Claude reviews every small change

**Solutions**:
- Add path filters to limit reviewed files
- Use conditional logic to skip minor changes
- Add `[skip-review]` to PR titles when needed

#### 4. Authentication Errors
**Problem**: "Permission denied" or authentication failures

**Solutions**:
```bash
# Refresh GitHub CLI authentication
gh auth refresh -h github.com -s repo,workflow

# Verify repository access
gh repo view Gaetan-Jaminon/GoTuto

# Check secret configuration
gh secret list
```

### Debugging Workflow Issues

**Check workflow runs**:
```bash
# List recent workflow runs
gh run list --workflow=claude.yml

# View specific run details
gh run view <run-id>

# Download logs
gh run download <run-id>
```

**Validate workflow syntax**:
```bash
# Use GitHub CLI to validate
gh workflow view claude.yml

# Check YAML syntax locally
yamllint .github/workflows/claude.yml
```

## 🎓 Best Practices

### 1. Effective @claude Mentions
- **Be specific**: Include relevant code snippets and error messages
- **Provide context**: Explain what you're trying to achieve
- **Ask focused questions**: Break complex problems into smaller parts

### 2. Custom Instructions
```yaml
custom_instructions: |
  Project Context:
  - Go microservices architecture
  - PostgreSQL database with migrations
  - OpenShift deployment
  - Enterprise security requirements
  
  Code Standards:
  - Follow Go naming conventions
  - Use structured logging
  - Implement comprehensive error handling
  - Write table-driven tests
  - Document all public APIs
  
  Review Focus:
  - Database query efficiency
  - Security vulnerabilities
  - Memory leaks and performance
  - Test coverage and quality
```

### 3. Review Quality
- Enable `use_sticky_comment: true` for persistent feedback
- Use file path filters for large repositories
- Customize prompts for different file types
- Set up different workflows for different environments

### 4. Team Collaboration
- Use Claude for initial reviews, human reviewers for final approval
- Tag `@claude` for quick questions during code reviews
- Create issues with `@claude` for architecture discussions
- Use Claude to onboard new team members

## 🔮 Future Enhancements

### Potential Improvements

#### 1. Advanced Workflow Integration
- **Automatic fix suggestions**: Claude creates fix commits
- **Test generation**: Auto-generate tests for new code
- **Documentation updates**: Keep docs in sync with code changes
- **Migration assistance**: Help with database schema changes

#### 2. Project-Specific Customization
- **Go-specific prompts**: Tailored for Go best practices
- **Database review focus**: PostgreSQL optimization suggestions
- **OpenShift deployment checks**: Container and security validation
- **API design guidance**: REST API best practices

#### 3. Enhanced Monitoring
- **Review quality metrics**: Track Claude's suggestion accuracy
- **Response time optimization**: Faster review processing
- **Integration analytics**: Usage patterns and effectiveness
- **Custom dashboards**: Visual insights into code quality trends

#### 4. Extended Capabilities
```yaml
# Future configuration options
enhanced_features:
  auto_fix_simple_issues: true
  generate_missing_tests: true
  update_documentation: true
  suggest_refactoring: true
  performance_profiling: true
```

### Roadmap Integration
The Claude integration will evolve alongside the broader CI/CD pipeline:

1. **Phase 1** (Current): Basic review and assistance
2. **Phase 2**: Integration with existing CI/CD workflows
3. **Phase 3**: Advanced automation and fix suggestions
4. **Phase 4**: Full AI-assisted development workflow

## 📚 Learning Resources

### Claude Code Documentation
- [Official Claude Code Docs](https://docs.anthropic.com/en/docs/claude-code)
- [GitHub Actions Integration](https://github.com/anthropics/claude-code-action)
- [Best Practices Guide](https://docs.anthropic.com/en/docs/claude-code/common-workflows)

### Go Development with AI
- Using AI for Go code reviews
- AI-assisted debugging patterns
- Test generation strategies
- Documentation automation

### GitHub Integration Patterns
- Advanced workflow triggers
- Secret management
- Permission configurations
- Multi-repository setups

## 🤝 Contributing to the Integration

### Improving Workflows
1. **Test changes locally** with workflow simulation tools
2. **Update documentation** when modifying configurations
3. **Share improvements** with the team
4. **Monitor effectiveness** and gather feedback

### Feedback Loop
- Track which suggestions are most helpful
- Identify common issues Claude catches
- Refine custom instructions based on experience
- Share learnings with the development team

---

This Claude/GitHub integration represents a significant enhancement to the development workflow, providing AI-powered assistance for code quality, debugging, and learning. As the integration matures, it will become an invaluable tool for maintaining high-quality Go code and accelerating development productivity.

**Next Steps**: Try mentioning `@claude` in a GitHub issue or PR to experience the integration firsthand!