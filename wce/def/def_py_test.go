package def

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/xackery/quail/helper"
	"gopkg.in/yaml.v3"
)

var (
	knownProps = make(map[string]bool)
)

func pyFloatExpr(expr string) string {
	return fmt.Sprintf("format(%s, '.8e')", expr)
}

func pyValueExpr(expr string, format string, nullable bool) string {
	switch format {

	case "%0.8e":
		if nullable {
			return fmt.Sprintf("('NULL' if %s is None else %s)", expr, pyFloatExpr(expr))
		}
		return pyFloatExpr(expr)

	default: // int / string fallback
		if nullable {
			return fmt.Sprintf("('NULL' if %s is None else %s)", expr, expr)
		}
		return expr
	}
}

func pyDefaultValue(prop Property) string {
	if len(prop.Properties) > 0 {
		return "" // handled elsewhere
	}

	if strings.Contains(pyPropType(prop), "list[") {
		if strings.HasSuffix(prop.Name, "?") {
			return "None"
		}
		return "[]"
	}

	isNullable := strings.HasSuffix(prop.Name, "?")

	// If nullable → default is None
	if isNullable {
		return "None"
	}

	// Otherwise derive from base type
	if len(prop.Args) == 1 {
		switch prop.Args[0].Format {
		case `%d`:
			return "0"
		case `%0.8e`:
			return "0.0"
		case `%s`:
			return "\"\""
		}
	}

	if len(prop.Args) > 1 {
		// tuple default
		parts := []string{}
		for _, arg := range prop.Args {
			switch arg.Format {
			case `%d`:
				parts = append(parts, "0")
			case `%0.8e`:
				parts = append(parts, "0.0")
			case `%s`:
				parts = append(parts, "\"\"")
			}
		}
		return fmt.Sprintf("(%s)", strings.Join(parts, ", "))
	}

	return "None"
}

func TestWceGenPython(t *testing.T) {
	// defs declared in def_md_test.go
	dirTest := helper.DirTest()

	for _, def := range defs {
		defName := strings.ToLower(def.Definition())

		r, err := os.Open(fmt.Sprintf("%s/../wce/def/%s.yaml", dirTest, strings.ToLower(defName)))
		if err != nil {
			t.Fatalf("open %s: %v", defName, err)
		}
		defer r.Close()

		yamlDef := &Definition{}
		err = yaml.NewDecoder(r).Decode(yamlDef)
		if err != nil {
			t.Fatalf("decode %s: %v", defName, err)
		}

		w, err := os.Create(fmt.Sprintf("%s/%s.py", dirTest, defName))
		if err != nil {
			t.Fatalf("create failed: %s", err)
		}
		defer w.Close()

		w.WriteString(`# Generated from quail, DO NOT EDIT
import io
from .parse import property

class ` + defName + `:
	@staticmethod
	def definition():
		return "` + strings.ToUpper(defName) + `"

`)
		if yamlDef.HasTag {
			w.WriteString("\ttag:str\n")
		}

		w.WriteString(initPyGen(yamlDef))

		decInitBuf := &bytes.Buffer{}
		propInitBuf := &bytes.Buffer{}
		initReaderBuf := &bytes.Buffer{}
		initWriterBuf := &bytes.Buffer{}
		initWriterBuf.WriteString("\tdef read(self, tag:str, r:io.TextIOWrapper|None) -> str:\n")
		if yamlDef.HasTag {
			initWriterBuf.WriteString("\t\tself.tag = tag\n")
		}

		writeBuf := &bytes.Buffer{}
		writeBuf.WriteString("\tdef write(self, w:io.TextIOWrapper)->str:\n")
		if yamlDef.HasTag {
			writeBuf.WriteString("\t\tw.write(f\"{self.definition()} \\\"{self.tag}\\\"\\n\")\n")
		} else {
			writeBuf.WriteString("\t\tw.write(f\"{self.definition()}\\n\")\n")
		}
		err = wcePyGen(propInitBuf, decInitBuf, initReaderBuf, initWriterBuf, writeBuf, yamlDef)
		if err != nil {
			t.Fatalf("wceGen %s: %v", defName, err)
		}
		initWriterBuf.WriteString("\t\tif r is None:\n")
		initWriterBuf.WriteString("\t\t\treturn \"no reader provided\"\n")

		w.Write(initWriterBuf.Bytes())
		w.WriteString("\n")
		w.Write(initReaderBuf.Bytes())
		w.WriteString("\t\treturn \"\"\n")
		w.WriteString("\n")
		w.Write(writeBuf.Bytes())
		w.WriteString("\t\treturn \"\"\n")
		w.WriteString("\n")
	}

	fmt.Println("Generated", len(defs), "definitions")

}

func initPyGen(yamlDef *Definition) string {
	out := ""

	out += initPyProperties(2, yamlDef.Properties, "self", yamlDef.HasTag)
	out += "\n"
	return out
}

func initPyProperties(tabIndex int, props []Property, scope string, hasTag bool) string {
	out := ""
	tabIndex--
	for _, prop := range props {
		if len(prop.Properties) == 0 && len(prop.Args) == 0 {
			continue
		}

		if len(prop.Properties) > 0 {
			continue
		}
		propName := strings.TrimSuffix(prop.Name, "?")
		propName = strings.ToLower(propName)
		propType := pyPropType(prop)
		if propType == "None" {
			continue
		}
		out += fmt.Sprintf("%s%s:%s\n", strings.Repeat("\t", tabIndex), propName, propType)
	}
	out += "\n"

	tabIndex++
	out += fmt.Sprintf("%sdef __init__(self):\n", strings.Repeat("\t", tabIndex-1))
	if hasTag {
		out += fmt.Sprintf("%sself.tag = \"\"\n", strings.Repeat("\t", tabIndex))
	}

	for _, prop := range props {
		if len(prop.Properties) == 0 && len(prop.Args) == 0 {
			continue
		}

		if len(prop.Properties) > 0 {
			continue
		}
		propName := strings.TrimSuffix(prop.Name, "?")
		propName = strings.ToLower(propName)
		propValue := pyDefaultValue(prop)
		if propValue == "" {
			continue
		}

		out += fmt.Sprintf("%sself.%s = %s #%d\n", strings.Repeat("\t", tabIndex), propName, propValue, tabIndex)
	}

	for _, prop := range props {
		if len(prop.Properties) == 0 {
			continue
		}

		propName := strings.TrimSuffix(prop.Name, "?")
		propNameLower := strings.ToLower(propName)

		isArray := len(prop.Args) == 1 &&
			prop.Args[0].Format == "%d"

		isSection := len(prop.Args) == 0

		if isArray {
			if strings.HasPrefix(propNameLower, "num") {
				propNameLower = strings.TrimPrefix(propNameLower, "num")
			}

			out += fmt.Sprintf("%sself.%s = []\n",
				strings.Repeat("\t", tabIndex),
				propNameLower,
			)
		}

		if isSection {
			out += fmt.Sprintf("%sself.%s = self.%s()\n",
				strings.Repeat("\t", tabIndex),
				propNameLower,
				propNameLower,
			)
		}
	}

	for _, prop := range props {
		if len(prop.Properties) == 0 {
			continue
		}

		propName := strings.TrimSuffix(prop.Name, "?")
		propNameLower := strings.ToLower(propName)

		isArray := len(prop.Args) == 1 &&
			prop.Args[0].Format == "%d"

		if isArray {

			elementName := strings.ToLower(prop.Properties[0].Name)

			out += "\n"
			out += fmt.Sprintf("%sclass %s:\n",
				strings.Repeat("\t", tabIndex-1),
				elementName,
			)

			out += initPyProperties(tabIndex+1, prop.Properties, elementName, false)

			continue
		}

		// normal section
		out += "\n"
		out += fmt.Sprintf("%sclass %s:\n",
			strings.Repeat("\t", tabIndex-1),
			propNameLower,
		)

		out += initPyProperties(tabIndex+1, prop.Properties, propNameLower, false)
	}

	return out
}

func wcePyGen(propInitBuf *bytes.Buffer, decInitBuf *bytes.Buffer, initReaderBuf *bytes.Buffer, initWriterBuf *bytes.Buffer, writeBuf *bytes.Buffer, yamlDef *Definition) error {
	knownProps = make(map[string]bool)
	for i, prop := range yamlDef.Properties {

		if i == 0 {
			decInitBuf.WriteString(fmt.Sprintf("\tdef __init__(self):\n"))
		}
		err := traversePyProp(propInitBuf, decInitBuf, initReaderBuf, initWriterBuf, writeBuf, prop, "self", 2, 1, "")
		if err != nil {
			return err
		}
		if i == len(yamlDef.Properties)-1 {
			propInitBuf.Write(decInitBuf.Bytes())
			decInitBuf.Reset()
		}
	}

	return nil
}

func traversePyProp(propInitBuf *bytes.Buffer, decInitBuf *bytes.Buffer, initReaderBuf *bytes.Buffer, initWriterBuf *bytes.Buffer, writeBuf *bytes.Buffer, prop Property, scope string, initTabCount int, decTabCount int, treeScope string) error {

	if knownProps[prop.Name] {
		return fmt.Errorf("duplicate property: %s", prop.Name)
	}
	knownProps[prop.Name] = true
	propBuf := ""
	initBuf := ""

	propKey := strings.TrimSuffix(prop.Name, "?")
	propKey = strings.ToLower(propKey)
	if strings.HasPrefix(propKey, "num") {
		propKey = strings.TrimPrefix(propKey, "num")
	}

	if treeScope == "" {
		treeScope = "self."
	}
	treeScope += propKey[:len(propKey)-1] + "."

	isNullable := strings.HasSuffix(prop.Name, "?")
	trimName := strings.TrimSuffix(prop.Name, "?")

	isArray := len(prop.Properties) > 0 &&
		len(prop.Args) == 1 &&
		prop.Args[0].Format == "%d"

	isSection := len(prop.Properties) > 0 &&
		len(prop.Args) == 0

	isManyArg := false
	if len(prop.Args) > 0 {
		if len(prop.Properties) == 0 {
			propBuf += strings.Repeat("\t", decTabCount) + strings.ToLower(trimName) + ":"
			initBuf += fmt.Sprintf("%s\tself.%s = ", strings.Repeat("\t", decTabCount), strings.ToLower(trimName))
		}
		if len(prop.Args) > 1 {
			propBuf += "tuple["
			initBuf += "tuple["
		}

		for _, arg := range prop.Args {

			if strings.HasSuffix(arg.Format, "...") {
				isManyArg = true
				arg.Format = strings.TrimSuffix(arg.Format, "...")
			}
		}

		argLen := len(prop.Args)

		if !isManyArg {
			initReaderBuf.WriteString(fmt.Sprintf("%srecords = property(r, \"%s\", %d)\n", strings.Repeat("\t", initTabCount), prop.Name, argLen))
			if len(prop.Properties) == 0 {
				initReaderBuf.WriteString(fmt.Sprintf("%s%s.%s = ", strings.Repeat("\t", initTabCount), scope, strings.ToLower(trimName)))
				propVar := fmt.Sprintf("%s.%s", scope, strings.ToLower(trimName))
				if len(prop.Args) > 1 {
					if isNullable {
						parts := []string{}
						for i := range prop.Args {
							expr := fmt.Sprintf("%s[%d]", propVar, i)
							formatted := pyValueExpr(expr, prop.Args[i].Format, false) // IMPORTANT: false here
							parts = append(parts, fmt.Sprintf("{('NULL' if %s is None else %s)}", propVar, formatted))
						}

						writeBuf.WriteString(fmt.Sprintf(
							"%sw.write(f\"%s%s %s\\n\")\n",
							strings.Repeat("\t", initTabCount),
							strings.Repeat("\\t", decTabCount),
							prop.Name,
							strings.Join(parts, " "),
						))
					} else {
						parts := []string{}
						for i := range prop.Args {
							expr := fmt.Sprintf("%s[%d]", propVar, i)
							parts = append(parts, fmt.Sprintf("{%s}", pyValueExpr(expr, prop.Args[i].Format, isNullable)))
						}

						writeBuf.WriteString(fmt.Sprintf(
							"%sw.write(f\"%s%s %s\\n\")\n",
							strings.Repeat("\t", initTabCount),
							strings.Repeat("\\t", decTabCount),
							prop.Name,
							strings.Join(parts, " "),
						))
					}
				} else {
					// single value → keep old behavior (for now)
					argFormat := prop.Args[0].Format
					propVar := fmt.Sprintf("%s.%s", scope, strings.ToLower(trimName))

					if argFormat == "%s" {
						if isNullable {

							writeBuf.WriteString(fmt.Sprintf(
								"%sif %s is None: w.write(\"%s%s NULL\\n\")\n",
								strings.Repeat("\t", initTabCount),
								propVar,
								strings.Repeat("\\t", decTabCount),
								prop.Name,
							))

							writeBuf.WriteString(fmt.Sprintf(
								"%selse: w.write(f\"%s%s \\\"{%s}\\\"\\n\")\n",
								strings.Repeat("\t", initTabCount),
								strings.Repeat("\\t", decTabCount),
								prop.Name,
								propVar,
							))

						} else {

							writeBuf.WriteString(fmt.Sprintf(
								"%sw.write(f\"%s%s \\\"{%s}\\\"\\n\")\n",
								strings.Repeat("\t", initTabCount),
								strings.Repeat("\\t", decTabCount),
								prop.Name,
								propVar,
							))
						}
					} else {
						// int / float → no quotes
						expr := pyValueExpr(propVar, argFormat, isNullable)

						writeBuf.WriteString(fmt.Sprintf(
							"%sw.write(f\"%s%s {%s}\\n\")\n",
							strings.Repeat("\t", initTabCount),
							strings.Repeat("\\t", decTabCount),
							prop.Name,
							expr,
						))
					}
				}
			} else {
				initReaderBuf.WriteString(fmt.Sprintf("%s%s = ", strings.Repeat("\t", initTabCount), strings.ToLower(trimName)))
			}
			for i, arg := range prop.Args {

				if strings.HasSuffix(arg.Format, "...") {
					isManyArg = true
					arg.Format = strings.TrimSuffix(arg.Format, "...")
				}

				base := ""
				switch arg.Format {
				case `%s`:
					base = "str"
				case `%d`:
					base = "int"
				case `%0.8e`:
					base = "float"
				default:
					return fmt.Errorf("unhandled type: %s", arg.Format)
				}

				if len(prop.Properties) == 0 {

					// -------------------------
					// SINGLE VALUE (keep old behavior)
					// -------------------------
					if len(prop.Args) == 1 {

						if isNullable {
							propBuf += fmt.Sprintf("%s | None", base)
							initBuf += fmt.Sprintf("%s | None", base)

							initReaderBuf.WriteString(fmt.Sprintf(
								"%s(records[%d]) if records[%d] != \"NULL\" else None",
								base, i+1, i+1))

						} else {
							propBuf += base

							switch base {
							case "int":
								initBuf += "0"
							case "float":
								initBuf += "0.0"
							case "str":
								initBuf += "\"\""
							}

							initReaderBuf.WriteString(fmt.Sprintf("%s(records[%d])", base, i+1))
						}

						// -------------------------
						// TUPLE (THIS IS THE FIX)
						// -------------------------
					} else {

						// Build correct tuple type ONCE (only when i == 0)
						if i == 0 {

							types := []string{}
							for _, a := range prop.Args {
								switch a.Format {
								case `%d`:
									types = append(types, "int")
								case `%0.8e`:
									types = append(types, "float")
								case `%s`:
									types = append(types, "str")
								}
							}

							tupleType := fmt.Sprintf("tuple[%s]", strings.Join(types, ", "))

							if isNullable {
								tupleType += " | None"
							}

							propBuf += tupleType
							initBuf += tupleType

							// ---- READER (ALL-OR-NOTHING) ----
							if isNullable {
								initReaderBuf.WriteString("None if records[1] == \"NULL\" else (")
							} else {
								initReaderBuf.WriteString("(")
							}
						}

						// Write each tuple element
						initReaderBuf.WriteString(fmt.Sprintf("%s(records[%d])", base, i+1))

						if i < len(prop.Args)-1 {
							initReaderBuf.WriteString(", ")
						} else {
							initReaderBuf.WriteString(")")
						}
					}
				} else {
					if isNullable {
						initReaderBuf.WriteString(fmt.Sprintf(
							"%s(records[%d]) if records[%d] != \"NULL\" else None\n",
							base,
							i+1,
							i+1,
						))
					} else {
						initReaderBuf.WriteString(fmt.Sprintf(
							"%s(records[%d])\n",
							base,
							i+1,
						))
					}
				}
				if len(prop.Args) > i+1 {

					// Only append to type strings for SINGLE VALUE mode
					if len(prop.Args) == 1 {
						propBuf += ", "
						initBuf += ", "
						initReaderBuf.WriteString(", ")
					}
				}
			}
		} else { // many args
			initReaderBuf.WriteString(fmt.Sprintf("%srecords = property(r, \"%s\", -1)\n", strings.Repeat("\t", initTabCount), prop.Name))
			if len(prop.Properties) == 0 {
				initReaderBuf.WriteString(fmt.Sprintf("%s%s.%s = ", strings.Repeat("\t", initTabCount), scope, strings.ToLower(trimName)))
				propVar := fmt.Sprintf("%s.%s", scope, strings.ToLower(trimName))
				if len(prop.Args) > 1 {
					if isNullable {
						parts := []string{}
						for i := range prop.Args {
							expr := fmt.Sprintf("%s[%d]", propVar, i)
							formatted := pyValueExpr(expr, prop.Args[i].Format, false) // IMPORTANT: false here
							parts = append(parts, fmt.Sprintf("{('NULL' if %s is None else %s)}", propVar, formatted))
						}

						writeBuf.WriteString(fmt.Sprintf(
							"%sw.write(f\"%s%s %s\\n\")\n",
							strings.Repeat("\t", initTabCount),
							strings.Repeat("\\t", decTabCount),
							prop.Name,
							strings.Join(parts, " "),
						))
					} else {
						parts := []string{}
						for i := range prop.Args {
							expr := fmt.Sprintf("%s[%d]", propVar, i)
							parts = append(parts, fmt.Sprintf("{%s}", pyValueExpr(expr, prop.Args[i].Format, isNullable)))
						}

						writeBuf.WriteString(fmt.Sprintf(
							"%sw.write(f\"%s%s %s\\n\")\n",
							strings.Repeat("\t", initTabCount),
							strings.Repeat("\\t", decTabCount),
							prop.Name,
							strings.Join(parts, " "),
						))
					}
				} else {
					argFormat := prop.Args[0].Format
					propVar := fmt.Sprintf("%s.%s", scope, strings.ToLower(trimName))

					if strings.HasSuffix(argFormat, "...") {

						if isNullable {

							writeBuf.WriteString(fmt.Sprintf(
								"%sw.write(f\"%s%s {'NULL' if %s is None else ' '.join(%s)}\\n\")\n",
								strings.Repeat("\t", initTabCount),
								strings.Repeat("\\t", decTabCount),
								prop.Name,
								propVar,
								propVar,
							))

						} else {

							writeBuf.WriteString(fmt.Sprintf(
								"%sw.write(f\"%s%s {' '.join(%s)}\\n\")\n",
								strings.Repeat("\t", initTabCount),
								strings.Repeat("\\t", decTabCount),
								prop.Name,
								propVar,
							))
						}
					} else if argFormat == "%s" {

						if isNullable {

							writeBuf.WriteString(fmt.Sprintf(
								"%sif %s is None: w.write(\"%s%s NULL\\n\")\n",
								strings.Repeat("\t", initTabCount),
								propVar,
								strings.Repeat("\\t", decTabCount),
								prop.Name,
							))

							writeBuf.WriteString(fmt.Sprintf(
								"%selse: w.write(f\"%s%s \\\"{%s}\\\"\\n\")\n",
								strings.Repeat("\t", initTabCount),
								strings.Repeat("\\t", decTabCount),
								prop.Name,
								propVar,
							))

						} else {

							writeBuf.WriteString(fmt.Sprintf(
								"%sw.write(f\"%s%s \\\"{%s}\\\"\\n\")\n",
								strings.Repeat("\t", initTabCount),
								strings.Repeat("\\t", decTabCount),
								prop.Name,
								propVar,
							))
						}

					} else {
						// int / float → no quotes (existing behavior)
						expr := pyValueExpr(propVar, argFormat, isNullable)

						writeBuf.WriteString(fmt.Sprintf(
							"%sw.write(f\"%s%s {%s}\\n\")\n",
							strings.Repeat("\t", initTabCount),
							strings.Repeat("\\t", decTabCount),
							prop.Name,
							expr,
						))
					}
				}
			} else {
				initReaderBuf.WriteString(fmt.Sprintf("%s%s = ", strings.Repeat("\t", initTabCount), strings.ToLower(trimName)))
			}
			if isNullable {
				propBuf += "list[str] | None"
				initBuf += "list[str] | None"
			} else {
				propBuf += "list[str]"
				initBuf += "list[str]"
			}
			if isNullable {
				initReaderBuf.WriteString(
					"None if len(records) > 1 and records[1] == \"NULL\" else records[1:]\n",
				)
			} else {
				initReaderBuf.WriteString("records[1:]\n")
			}
		}

		if prop.Note != "" && len(prop.Properties) == 0 {
			propBuf += " # " + prop.Note
		}
		propBuf += "\n"
		initBuf += "\n"
	} else { // no argument parse
		initReaderBuf.WriteString(fmt.Sprintf("%sproperty(r, \"%s\", 0)\n", strings.Repeat("\t", initTabCount), prop.Name))
		writeBuf.WriteString(fmt.Sprintf("%sw.write(f\"%s%s\\n\")\n", strings.Repeat("\t", initTabCount), strings.Repeat("\\t", decTabCount), prop.Name))
	}

	initReaderBuf.WriteString("\n")

	decInitBuf.WriteString(initBuf)
	propInitBuf.WriteString(propBuf)

	if isArray {

		lastScope := scope

		// container name from NUMXXXX (e.g. NUMFRAMES -> frames)
		containerName := strings.ToLower(trimName)
		if strings.HasPrefix(containerName, "num") {
			containerName = strings.TrimPrefix(containerName, "num")
		}

		// element name from first child property (e.g. FRAME)
		elementName := strings.ToLower(prop.Properties[0].Name)

		// __init__: create empty list
		decInitBuf.WriteString(fmt.Sprintf(
			"%sself.%s = []\n",
			strings.Repeat("\t", initTabCount),
			containerName,
		))

		// define nested element class at CLASS scope (not inside __init__)
		// propInitBuf.WriteString(fmt.Sprintf(
		// 	"\tclass %s:\n",
		// 	elementName,
		// ))
		// propInitBuf.WriteString(fmt.Sprintf(
		// 	"\t\tdef __init__(self):\n",
		// ))

		// reset list during read
		initReaderBuf.WriteString(fmt.Sprintf(
			"%s%s.%s = []\n",
			strings.Repeat("\t", initTabCount),
			lastScope,
			containerName,
		))

		// loop count (NUMXXXX already parsed into lowercase trimName)
		if isNullable {
			initReaderBuf.WriteString(fmt.Sprintf(
				"%sfor %s in range(%s or 0):\n",
				strings.Repeat("\t", initTabCount),
				tabCode(initTabCount),
				strings.ToLower(trimName),
			))
		} else {
			initReaderBuf.WriteString(fmt.Sprintf(
				"%sfor %s in range(%s):\n",
				strings.Repeat("\t", initTabCount),
				tabCode(initTabCount),
				strings.ToLower(trimName),
			))
		}

		// instantiate nested element properly
		parentVar := scope
		// if strings.Contains(parentVar, ".") {
		// 	parts := strings.Split(parentVar, ".")
		// 	parentVar = parts[len(parts)-1]
		// }

		initReaderBuf.WriteString(fmt.Sprintf(
			"%s\t%s = type(%s).%s()\n",
			strings.Repeat("\t", initTabCount),
			elementName+tabCode(initTabCount),
			parentVar,
			elementName,
		))

		// write count
		writeBuf.WriteString(fmt.Sprintf(
			"%sw.write(f\"%s%s {len(%s.%s)}\\n\")\n",
			strings.Repeat("\t", initTabCount),
			strings.Repeat("\\t", decTabCount),
			prop.Name,
			lastScope,
			containerName,
		))

		// write loop
		writeBuf.WriteString(fmt.Sprintf(
			"%sfor %s in %s.%s:\n",
			strings.Repeat("\t", initTabCount),
			elementName+tabCode(initTabCount),
			lastScope,
			containerName,
		))

		// recurse into child properties
		for _, prop2 := range prop.Properties {
			err := traversePyProp(
				propInitBuf,
				decInitBuf,
				initReaderBuf,
				initWriterBuf,
				writeBuf,
				prop2,
				elementName+tabCode(initTabCount),
				initTabCount+1,
				decTabCount+1,
				treeScope,
			)
			if err != nil {
				return err
			}
		}

		// append element to container
		initReaderBuf.WriteString(fmt.Sprintf(
			"%s\t%s.%s.append(%s)\n",
			strings.Repeat("\t", initTabCount),
			lastScope,
			containerName,
			elementName+tabCode(initTabCount),
		))
	}

	if isSection {

		sectionName := strings.ToLower(trimName)

		// create section instance in __init__
		decInitBuf.WriteString(fmt.Sprintf(
			"%sself.%s = self.%s()\n",
			strings.Repeat("\t", initTabCount),
			sectionName,
			trimName,
		))

		// define class
		decInitBuf.WriteString(fmt.Sprintf(
			"%sclass %s:\n",
			strings.Repeat("\t", decTabCount),
			trimName,
		))
		decInitBuf.WriteString(fmt.Sprintf(
			"%s\tdef __init__(self):\n",
			strings.Repeat("\t", decTabCount),
		))

		// read header
		// initReaderBuf.WriteString(fmt.Sprintf(
		// 	"%sproperty(r, \"%s\", 0)\n",
		// 	strings.Repeat("\t", initTabCount),
		// 	prop.Name,
		// ))

		// recurse WITHOUT extra indentation offset
		for _, prop2 := range prop.Properties {
			err := traversePyProp(
				propInitBuf,
				decInitBuf,
				initReaderBuf,
				initWriterBuf,
				writeBuf,
				prop2,
				scope+"."+sectionName,
				initTabCount,
				decTabCount+1,
				treeScope,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func tabCode(tabCount int) string {
	return string(rune('i' + tabCount - 2))
}

func pyPropType(prop Property) string {
	if len(prop.Properties) == 0 && len(prop.Args) == 0 {
		return "None"
	}

	// Handle variable-length args (list)
	for _, arg := range prop.Args {
		if strings.HasSuffix(arg.Format, "...") {
			if strings.HasSuffix(prop.Name, "?") {
				return "list[str] | None"
			}
			return "list[str]"
		}
	}

	isNullable := strings.HasSuffix(prop.Name, "?")

	// Build list of base types
	types := []string{}
	for _, arg := range prop.Args {
		switch arg.Format {
		case `%d`:
			types = append(types, "int")
		case `%0.8e`:
			types = append(types, "float")
		case `%s`:
			types = append(types, "str")
		default:
			types = append(types, "Unknown")
		}
	}

	out := ""

	// Single value
	if len(types) == 1 {
		out = types[0]
	} else {
		// Tuple
		out = fmt.Sprintf("tuple[%s]", strings.Join(types, ", "))
	}

	// Nullable
	if isNullable {
		out += " | None"
	}

	return out
}
