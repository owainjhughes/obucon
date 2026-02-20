#!/usr/bin/env python3
"""
Transform JLPT level CSV files into PostgreSQL migration format.
Generates migration SQL files compatible with golang-migrate.
"""

import csv
from pathlib import Path


def jlpt_level_to_int(level_str):
    mapping = {"n5": 5, "n4": 4, "n3": 3, "n2": 2, "n1": 1}
    return mapping.get(level_str.lower(), None)


def escape_sql_string(s):
    if s is None or s == "":
        return "NULL"
    escaped = s.replace("'", "''")
    return f"'{escaped}'"


def transform_jlpt_data():
    input_dir = Path(__file__).parent / "original_data"
    csv_output = Path(__file__).parent / "jlpt_vocabulary_seed.csv"
    migration_dir = Path(__file__).parent / ".." / "migrations"
    migration_dir.mkdir(parents=True, exist_ok=True)

    migration_up = migration_dir / "002_seed_jlpt_vocabulary.up.sql"
    migration_down = migration_dir / "002_seed_jlpt_vocabulary.down.sql"

    levels = ["n5", "n4", "n3", "n2", "n1"]
    csv_rows = []
    insert_statements = []

    for level in levels:
        input_file = input_dir / f"{level}.csv"
        print(f"Processing {level.upper()}: {input_file}")

        if not input_file.exists():
            print(f"  Warning: {input_file} not found, skipping...")
            continue

        jlpt_int = jlpt_level_to_int(level)
        with open(input_file, "r", encoding="utf-8") as f:
            reader = csv.DictReader(f)
            count = 0

            for row in reader:
                hiragana = row.get("kana", "").strip()
                kanji = row.get("kanji", "").strip()
                meaning = row.get("waller_definition", "").strip()

                kanji_form = kanji if kanji else hiragana

                if kanji_form and meaning:
                    csv_rows.append({
                        "kanji": kanji_form,
                        "hiragana": hiragana,
                        "meaning": meaning,
                        "jlpt_level": jlpt_int,
                    })
                    
                    insert_stmt = f"INSERT INTO japanese_dictionary (kanji, hiragana, meaning, jlpt_level) VALUES ({escape_sql_string(kanji_form)}, {escape_sql_string(hiragana)}, {escape_sql_string(meaning)}, {jlpt_int});"
                    insert_statements.append(insert_stmt)
                    count += 1

            print(f"  Added {count} entries")

    print(f"\nWriting CSV seed data to {csv_output}")
    with open(csv_output, "w", newline="", encoding="utf-8") as f:
        writer = csv.DictWriter(f, fieldnames=["kanji", "hiragana", "meaning", "jlpt_level"])
        writer.writeheader()
        for row in csv_rows:
            writer.writerow(row)

    print(f"Writing migration UP to {migration_up}")
    with open(migration_up, "w", encoding="utf-8") as f:
        f.write("Seed JLPT vocabulary into japanese_dictionary table\n")
        f.write("Compatible with golang-migrate\n\n")
        f.write("BEGIN;\n\n")
        for stmt in insert_statements:
            f.write(stmt + "\n")
        f.write("\nCOMMIT;\n")

    print(f"Writing migration DOWN to {migration_down}")
    with open(migration_down, "w", encoding="utf-8") as f:
        f.write("Rollback: Remove seeded JLPT vocabulary\n")
        f.write("Only deletes pre-populated dictionary entries (jlpt_level IS NOT NULL)\n\n")
        f.write("DELETE FROM japanese_dictionary WHERE jlpt_level IS NOT NULL;\n")

    print(f"\nTotal entries: {len(insert_statements)}")
    print("Done!")


if __name__ == "__main__":
    transform_jlpt_data()
