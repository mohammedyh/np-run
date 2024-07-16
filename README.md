# np-run

A CLI tool that lets you browse and execute npm, pnpm, yarn, or bun scripts from a local project.

https://github.com/mohammedyh/np-run/assets/32526267/4b866d0e-9e0a-4251-b874-ad1ab6d8df94

## Todo

- [x] Check for `package.json` in current directory / project
- [x] Read file and parse key and values in `scripts` object in `package.json`
- [x] Print scripts to stdout as a list - use BubbleTea
- [x] Detect the package manager - using the lockfile?
- [x] When enter is pressed for a chosen script, that script is run
- [ ] Look at how to properly structure a Go program

### Nice-To-Have Features
- [x] Install dependencies if `node_modules` isn't found
  - [x] Show spinner or progress bar while dependencies are installing
- [x] Print something if multiple lockfiles are found
