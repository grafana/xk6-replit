#!/usr/bin/env ts-node

import {
    Project,
    Node,
    InterfaceDeclaration,
    MethodSignature,
    PropertySignature,
    ConstructSignatureDeclaration,
    VariableDeclaration,
    ClassDeclaration,
    MethodDeclaration,
    PropertyDeclaration,
    ConstructorDeclaration,
} from "ts-morph";
import * as path from "path";
import * as fs from "fs";
import { parse as parseComment } from "comment-parser";

/* ----- JSDoc Helpers ----- */

/** Cleans raw JSDoc text by stripping /** markers and leading asterisks. */
function cleanJsDocText(text: string): string {
    return text
        .replace(/^\/\*\*[\r\n]*/, "")
        .replace(/\*\/$/, "")
        .replace(/^\s*\* ?/gm, "")
        .trim();
}

/** Removes {@link ...} markup, replacing it with its display text. */
function removeJsDocLinks(text: string): string {
    return text
        .replace(/\{@link\s+[^|]+\s*\|\s*([^}]+)\}/g, "$1")
        .replace(/\{@link\s+([^}]+)\}/g, "$1");
}

/** Parses a JSDoc block using comment‑parser. */
function parseJsDoc(jsDoc: any): { description: string; tags: any[] } | null {
    const raw = jsDoc.getFullText();
    const parsed = parseComment(raw);
    return parsed.length > 0 ? parsed[0] : null;
}

/** Returns the free‑form summary from a JSDoc block (text before any tags). */
function safeGetJsDocSummary(jsDoc: any): string {
    const parsed = parseJsDoc(jsDoc);
    if (!parsed) return "";
    return removeJsDocLinks(cleanJsDocText(parsed.description));
}

/** Returns a tag’s comment as a string. */
function safeGetTagComment(tag: any): string {
    const text = tag.description || "";
    return removeJsDocLinks(cleanJsDocText(text));
}

/** Returns an array of JSDoc nodes from a declaration, if available. */
function safeGetJsDocs(decl: any): any[] {
    if ("getJsDocs" in decl && typeof decl.getJsDocs === "function") {
        return decl.getJsDocs();
    }
    return [];
}

/** Simplifies a type string by removing import wrappers. */
function simplifyType(typeStr: string): string {
    return typeStr.replace(/import\(".*?"\)\./g, "");
}

/**
 * Extracts structured JSDoc info from an array of JSDoc nodes.
 */
function extractJsDocInfo(jsDocs: any[]): {
    summary: string;
    params: Record<string, string>;
    returns: string;
    example: string;
    deprecated: string;
    defaultValue: string;
} {
    let summary = "";
    const params: Record<string, string> = {};
    let returns = "";
    let example = "";
    let deprecated = "";
    let defaultValue = "";
    for (const jsDoc of jsDocs) {
        const parsed = parseJsDoc(jsDoc);
        if (!parsed) continue;
        summary += parsed.description.trim() + "\n";
        for (const tag of parsed.tags) {
            if (tag.tag === "param") {
                params[tag.name] = safeGetTagComment(tag);
            } else if (tag.tag === "returns" || tag.tag === "return") {
                returns = safeGetTagComment(tag);
            } else if (tag.tag === "example") {
                example = safeGetTagComment(tag);
            } else if (tag.tag === "deprecated") {
                deprecated = safeGetTagComment(tag);
            } else if (tag.tag === "defaultValue") {
                defaultValue = safeGetTagComment(tag);
            }
        }
    }
    return {
        summary: removeJsDocLinks(cleanJsDocText(summary.trim())),
        params,
        returns: removeJsDocLinks(returns),
        example: removeJsDocLinks(example),
        deprecated: removeJsDocLinks(deprecated),
        defaultValue: removeJsDocLinks(defaultValue),
    };
}

/* ----- Member Processing ----- */

/** Safely gets the name of a member (or returns overrideName if provided). */
function safeGetMemberName(member: any, overrideName?: string): string {
    if (overrideName) return overrideName;
    if (typeof member.getName === "function") return member.getName();
    if ("name" in member && typeof member.name === "string") return member.name;
    return "";
}

/**
 * Processes a member (method, property, or construct signature) generically.
 * The optional overrideName is used when processing a property from an object literal.
 */
function processMemberGeneric(member: any, overrideName?: string): any | null {
    if (Node.isMethodSignature(member)) {
        const method = member as MethodSignature;
        const name = safeGetMemberName(method, overrideName);
        const jsDocs = safeGetJsDocs(method);
        const jsDocInfo = extractJsDocInfo(jsDocs);
        const inputs = method.getParameters().map(param => {
            const paramName = param.getName();
            const paramType = simplifyType(param.getType().getText().replace(/\s+/g, " ").trim());
            return { name: paramName, type: paramType, documentation: jsDocInfo.params[paramName] || "" };
        });
        const signature = `${name}(${inputs.map(p => `${p.name}: ${p.type}`).join(", ")})`;
        const returnType = simplifyType(method.getReturnType().getText().replace(/\s+/g, " ").trim());
        return {
            kind: "method",
            name,
            signature,
            documentation: jsDocInfo.summary,
            params: inputs,
            returns: { type: returnType, documentation: jsDocInfo.returns },
        };
    } else if (Node.isPropertySignature(member)) {
        const prop = member as PropertySignature;
        const name = safeGetMemberName(prop, overrideName);
        const jsDocs = safeGetJsDocs(prop);
        const jsDocInfo = extractJsDocInfo(jsDocs);
        const typeText = simplifyType(prop.getType().getText().replace(/\s+/g, " ").trim());
        return {
            kind: "property",
            name,
            signature: prop.getText().replace(/\s+/g, " ").trim(),
            documentation: jsDocInfo.summary,
            type: typeText,
        };
    } else if (Node.isConstructSignatureDeclaration(member)) {
        // Process construct signatures generically as methods with overrideName (if provided) or "new"
        const name = overrideName || "new";
        const construct = member as ConstructSignatureDeclaration;
        const jsDocs = safeGetJsDocs(construct);
        const jsDocInfo = extractJsDocInfo(jsDocs);
        const inputs = construct.getParameters().map(param => {
            const paramName = param.getName();
            const paramType = simplifyType(param.getType().getText().replace(/\s+/g, " ").trim());
            return { name: paramName, type: paramType, documentation: jsDocInfo.params[paramName] || "" };
        });
        let typeParams = "";
        const tps = construct.getTypeParameters();
        if (tps.length > 0) {
            typeParams = "<" + tps.map(tp => tp.getText()).join(", ") + ">";
        }
        const signature = `${name}${typeParams}( ${inputs.map(p => `${p.name}: ${p.type}`).join(", ")} )`;
        const returnType = simplifyType(construct.getReturnType().getText().replace(/\s+/g, " ").trim());
        return {
            kind: "constructor",
            name,
            signature: `${signature}: ${returnType}`,
            documentation: jsDocInfo.summary,
            params: inputs,
            returns: { type: returnType, documentation: removeJsDocLinks(jsDocInfo.returns) },
        };
    }
    return null;
}

/* ----- Declaration Processing ----- */

/**
 * Processes an exported declaration generically.
 * - For interfaces, recurses through its members.
 * - For classes, processes its constructors and members.
 * - For variables whose type is an object literal, recurses through its properties.
 * - Otherwise, returns a simple declaration.
 */
function processDeclaration(decl: any, exportName: string): any {
    const kind = decl.getKindName();
    if (kind === "ClassDeclaration") {
        const cls = decl as ClassDeclaration;
        const result: any = { kind: "class", name: exportName, documentation: "", constructors: [], members: {} };
        const jsDocs = safeGetJsDocs(cls);
        const jsDocInfo = extractJsDocInfo(jsDocs);
        result.documentation = jsDocInfo.summary;
        result.signature = exportName;
        cls.getConstructors().forEach((ctor: ConstructorDeclaration) => {
            result.constructors.push(processConstructor(ctor));
        });
        cls.getMembers().forEach(member => {
            if (Node.isMethodDeclaration(member)) {
                result.members[member.getName()] = processMethod(member);
            } else if (Node.isPropertyDeclaration(member)) {
                result.members[member.getName()] = processProperty(member);
            }
        });
        if (cls.getExtends()) {
            result.extends = cls.getExtends()!.getText();
        }
        return result;
    } else if (kind === "InterfaceDeclaration") {
        const iface = decl as InterfaceDeclaration;
        const result: any = { kind: "interface", name: exportName, documentation: "", members: {} };
        const jsDocs = safeGetJsDocs(iface);
        const jsDocInfo = extractJsDocInfo(jsDocs);
        result.documentation = jsDocInfo.summary;
        result.signature = exportName;
        iface.getMembers().forEach(member => {
            if (
                Node.isMethodSignature(member) ||
                Node.isPropertySignature(member) ||
                Node.isConstructSignatureDeclaration(member)
            ) {
                const m = processMemberGeneric(member);
                if (!m) return;
                if (Node.isConstructSignatureDeclaration(member)) {
                    if (!result.members["new"]) result.members["new"] = [];
                    result.members["new"].push(m);
                } else {
                    result.members[member.getName()] = m;
                }
            }
        });
        return result;
    } else if (kind === "VariableDeclaration") {
        const varDecl = decl as VariableDeclaration;
        const result: any = { kind: "variable", name: exportName, documentation: "", signature: exportName };
        const jsDocs = safeGetJsDocs(varDecl);
        const jsDocInfo = extractJsDocInfo(jsDocs);
        result.documentation = jsDocInfo.summary;
        result.type = simplifyType(varDecl.getType().getText().replace(/\s+/g, " ").trim());
        const varType = varDecl.getType();
        if (varType.isObject() && varType.getProperties().length > 0) {
            result.members = {};
            for (const prop of varType.getProperties()) {
                const propDecls = prop.getDeclarations();
                const processed = propDecls.map(d => processMemberGeneric(d, prop.getName())).filter(x => x !== null);
                result.members[prop.getName()] = processed.length === 1 ? processed[0] : processed;
            }
        }
        return result;
    } else {
        // Generic declaration.
        let signature = simplifyType(decl.getText().replace(/\s+/g, " ").trim());
        signature = signature.replace(/^export\s+/, "");
        const jsDocs = safeGetJsDocs(decl);
        const jsDocInfo = extractJsDocInfo(jsDocs);
        return { kind: "declaration", name: exportName, signature, documentation: jsDocInfo.summary };
    }
}

/** Processes a constructor (for a class). */
function processConstructor(ctor: ConstructorDeclaration): any {
    const jsDocs = safeGetJsDocs(ctor);
    const jsDocInfo = extractJsDocInfo(jsDocs);
    const params = ctor.getParameters().map(param => {
        const paramName = param.getName();
        const paramType = simplifyType(param.getType().getText().replace(/\s+/g, " ").trim());
        return { name: paramName, type: paramType, documentation: jsDocInfo.params[paramName] || "" };
    });
    const signature = `constructor(${params.map(p => `${p.name}: ${p.type}`).join(", ")})`;
    return { kind: "constructor", signature, documentation: jsDocInfo.summary, params };
}

/** Processes a class method (MethodDeclaration). */
function processMethod(method: MethodDeclaration): any {
    // For class methods, use similar logic as for interface methods.
    const name = method.getName();
    const jsDocs = safeGetJsDocs(method);
    const jsDocInfo = extractJsDocInfo(jsDocs);
    const params = method.getParameters().map(param => {
        const paramName = param.getName();
        const paramType = simplifyType(param.getType().getText().replace(/\s+/g, " ").trim());
        return { name: paramName, type: paramType, documentation: jsDocInfo.params[paramName] || "" };
    });
    const signature = `${name}(${params.map(p => `${p.name}: ${p.type}`).join(", ")})`;
    const returnType = simplifyType(method.getReturnType().getText().replace(/\s+/g, " ").trim());
    return {
        kind: "method",
        name,
        signature,
        documentation: jsDocInfo.summary,
        params,
        returns: { type: returnType, documentation: jsDocInfo.returns }
    };
}

/** Processes a class property (PropertyDeclaration). */
function processProperty(prop: PropertyDeclaration): any {
    const name = prop.getName();
    const jsDocs = safeGetJsDocs(prop);
    const jsDocInfo = extractJsDocInfo(jsDocs);
    const typeText = simplifyType(prop.getType().getText().replace(/\s+/g, " ").trim());
    return {
        kind: "property",
        name,
        signature: prop.getText().replace(/\s+/g, " ").trim(),
        documentation: jsDocInfo.summary,
        type: typeText
    };
}

/**
 * Processes all exported declarations in a source file using getExportSymbols().
 */
function extractDocsHierarchical(sourceFilePath: string): any {
    const project = new Project({ compilerOptions: { allowJs: true } });
    const sourceFile = project.addSourceFileAtPath(sourceFilePath);
    const output: any = {};
    const exportSymbols = sourceFile.getExportSymbols();
    exportSymbols.forEach(symbol => {
        const exportName = symbol.getName();
        const decls = symbol.getDeclarations();
        const processed = decls.map(decl => processDeclaration(decl, exportName));
        // Merge multiple declarations if necessary.
        output[exportName] = processed.reduce((acc, cur) => ({ ...acc, ...cur }), {});
    });
    return output;
}

/**
 * Recursively finds all .d.ts files under a directory (skipping node_modules).
 */
function findDTSFiles(dir: string): string[] {
    let results: string[] = [];
    const list = fs.readdirSync(dir);
    list.forEach(file => {
        if (file === "node_modules") return;
        const filePath = path.join(dir, file);
        let stat;
        try {
            stat = fs.statSync(filePath);
        } catch (e) {
            console.error(`Error reading ${filePath}: ${e}`);
            return;
        }
        if (stat.isDirectory()) {
            results = results.concat(findDTSFiles(filePath));
        } else if (filePath.endsWith(".d.ts")) {
            results.push(filePath);
        }
    });
    return results;
}

/**
 * Returns a replacer for JSON.stringify that omits circular references.
 */
function getCircularReplacer() {
    const seen = new WeakSet();
    return function (key: string, value: any) {
        if (typeof value === "object" && value !== null) {
            if (seen.has(value)) return;
            seen.add(value);
        }
        return value;
    };
}

function main() {
    if (process.argv.length < 3) {
        console.error("Usage: ts-node parse-ts-defs.ts <input directory> [output directory]");
        process.exit(1);
    }
    const inputDir = path.resolve(process.argv[2]);
    const outputDir = process.argv[3] ? path.resolve(process.argv[3]) : inputDir;
    const dtsFiles = findDTSFiles(inputDir);
    console.log(`Found ${dtsFiles.length} .d.ts file(s).`);
    dtsFiles.forEach(filePath => {
        const relativePath = path.relative(inputDir, filePath).replace(/\\/g, "/");
        const outputRelativePath = relativePath.replace(/\.d\.ts$/, ".json");
        const outputFilePath = path.join(outputDir, outputRelativePath);
        fs.mkdirSync(path.dirname(outputFilePath), { recursive: true });
        const docs = extractDocsHierarchical(filePath);
        fs.writeFileSync(outputFilePath, JSON.stringify(docs, getCircularReplacer(), 2), "utf-8");
        console.log(`Wrote ${outputFilePath}`);
    });
}

main();
