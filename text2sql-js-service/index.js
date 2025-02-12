const express = require("express");
const { HfInference } = require("@huggingface/inference");
const { OpenAI } = require("openai");
const bodyParser = require("body-parser");
const dotenv = require("dotenv");
const path = require('path');
const localModelPath = path.join(__dirname, '../TextToSQL-Service/text-to-sql-simple');

dotenv.config();

const app = express();
app.use(bodyParser.json());

const hf = new HfInference(process.env.HUGGINGFACE_API_KEY);

const openai = new OpenAI({
    apiKey: process.env.OPENAI_API_KEY,
});

// Simple endpoint for T5 model
app.post("/text-to-sql/simple", async (req, res) => {
    const { query } = req.body;

    if (!query) {
        return res.status(400).json({ error: "Query cannot be empty" });
    }

    try {
        const inputText = `Translate English to SQL: ${query}`;
        const result = await hf.textGeneration({
            model: localModelPath,
            inputs: inputText,
            options: { wait_for_model: true },
        });

        if (!result.generated_text) {
            return res.status(500).json({ error: "Failed to generate SQL query" });
        }

        res.json({ sql_query: result.generated_text.replace(/"/g, "'") });
    } catch (error) {
        console.error("Error generating SQL:", error);
        res.status(500).json({ error: "Internal server error" });
    }
});

// Complex endpoint for T5 model
app.post("/text-to-sql/complex", async (req, res) => {
    const { schema, query } = req.body;

    if (!schema || !query) {
        return res.status(400).json({ error: "Schema and query cannot be empty" });
    }

    try {
        const schemaText = schema.tables
            .map((table) => {
                const columns = table.columns.map((col) => `${col.name} ${col.type}`).join(", ");
                const constraints = [];
                if (table.primaryKey) constraints.push(`PRIMARY KEY (${table.primaryKey})`);
                table.columns.forEach((col) => {
                    if (col.isForeignKey) {
                        constraints.push(
                            `FOREIGN KEY (${col.name}) REFERENCES ${col.referencedTable}(${col.referencedColumn})`
                        );
                    }
                });
                return `CREATE TABLE ${table.name} (${columns}${constraints.length > 0 ? `, ${constraints.join(", ")}` : ""});`;
            })
            .join(" ");

        const inputText = `Schema: ${schemaText}. Query for: ${query}`;
        const result = await hf.textGeneration({
            model: "text-to-sql-complex",
            inputs: inputText,
            options: { wait_for_model: true },
        });

        res.json({ sql_query: result.generated_text.replace(/"/g, "'") });
    } catch (error) {
        console.error("Error generating SQL:", error);
        res.status(500).json({ error: "Internal server error" });
    }
});

// GPT-3.5 Endpoint
app.post("/text-to-sql/gpt", async (req, res) => {
    const { schema, query } = req.body;

    if (!schema || !query) {
        return res.status(400).json({ error: "Schema and query cannot be empty" });
    }

    try {
        const schemaText = schema.tables
            .map(
                (table) =>
                    `Table ${table.name} has columns ${table.columns.map((col) => col.name).join(", ")}.`
            )
            .join(" ");

        const prompt = `Schema: ${schemaText}. Translate English to SQL: ${query}`;

        const response = await openai.createCompletion({
            model: "gpt-3.5-turbo",
            prompt,
            max_tokens: 2048,
            temperature: 1,
        });

        const sqlQuery = response.data.choices[0]?.text.trim();

        if (!sqlQuery) {
            return res.status(500).json({ error: "Failed to generate SQL query" });
        }

        res.json({ sql_query: sqlQuery });
    } catch (error) {
        console.error("Error generating SQL with GPT:", error);
        res.status(500).json({ error: "Internal server error" });
    }
});

app.get("/health", (req, res) => {
    res.json({ status: "healthy" });
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
    console.log(`Server running on http://localhost:${PORT}`);
});
