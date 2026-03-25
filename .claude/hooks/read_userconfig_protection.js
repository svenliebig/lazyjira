const fs = require("fs");
const path = require("path");
const log = path.join(__dirname, "log.txt");

async function main() {
  const chunks = [];

  for await (const chunk of process.stdin) {
    chunks.push(chunk);
  }

  const toolArgs = JSON.parse(Buffer.concat(chunks).toString());

  fs.appendFileSync(log, `reading: ${JSON.stringify(toolArgs)}\n`);

  // Extract the file path Claude is trying to read
  const readPath =
    toolArgs.tool_input?.file_path || toolArgs.tool_input?.path || "";

  if (
    readPath.includes("~/.config/lazyjira/config.json") ||
    readPath.includes(".config/lazyjira/config.json")
  ) {
    console.error(
      "You cannot read the config file, it contains sensitive information like the Jira API token.",
    );

    process.exit(2);
  }
}

main();
