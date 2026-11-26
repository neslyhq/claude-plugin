# Fixture: the build failed, the platform repaired it in its own build copy.

User: Deploy https://github.com/acme/shop to nesly.

Context for the agent under test:
- Signed in. `nesly deploy https://github.com/acme/shop --preview --detach --yes --json`
  returned a deployment id; `nesly wait <id> --json` exits 0 with status `ready`.
- The deploy log contains: `build failed: … the platform is diagnosing it`, then
  `platform fix: the start script pointed at a file that does not exist — changed
  package.json in the platform's build copy; your repository is unchanged`.
- `nesly fix show` prints a three-line diff of `package.json`.
- The repository has a `package-lock.json`; its `express` version is old.
