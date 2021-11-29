import { parseTemplate } from "@angular/compiler";
import fs from "fs";

const code = fs.readFileSync("../bin/template.html", "utf8");

const start = Date.now();

const data = parseTemplate(code, "./template.html");

const end = Date.now();

console.log(end - start);

// fs.writeFileSync("out.json", JSON.stringify(data));
