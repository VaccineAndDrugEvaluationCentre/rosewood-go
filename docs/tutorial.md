# Carpenter
Version 0.5.0

## Introduction
Rosewood is a domain-specific simple language primarily intended to simplify the automatic generation and formatting of statistical tables. The language provides facilities for representing tables in a simple, human-readable format and for manipulating and formatting the structure and style of those tables. Because Rosewood's source files are plain text files, they can be created by hand or automatically generated using a statistical package (e.g., SAS, R, Stata) or other general-purpose scripting or programming languages. 

Rosewood uses a simple markup, inspired by Markdown, to define a table and its contents making it easier to create and maintain than more complicated markup languages such as HTML or XML. The simple non-intrusive markup makes Rosewood tables easily readable even without rendering the source code or using special viewers. This feature was the overriding design goal of Rosewood as analysts need to frequently inspect the tables for correctness while analyzing the data, and forcing them to use a different tool to view the table contents tends to interrupt their work flow. For this reason, the code needed to format (e.g., merge certain cells) and style the table and its content is stored in a different section of the file, so not to reduce the readability of the tabular data. 

Carpenter is a cross-platform (Windows, MacOS, Linux) tool that can be used to parse and render Rosewood source files into html or docx files for viewing or printing.

## Create your first Rosewood file
- Rosewood files are text files, so they can be created using any text editor, e.g., Notepad (PC) or TextEdit (Mac). 
- Create a new text file in your favourite editor and write or copy and paste the following text into it.
```
+++
My first Rosewood table
+++
Item             | Stars |
Butter chicken   |   5   |
Star-anise candy |   4   |
Wilted lettuce   |   0   |
+++
footnote 1: butter chicken wins!
+++
+++
```
- Save the file in a folder (e.g., c:\Desktop), and name it simple.rw.
- Open a Command window or terminal and change to the folder where you saved the file.
- At the command prompt, type 
```
 ./carpenter run ../docs/simple.rw
```
OR 
```
 carpenter run ../docs/simple.rw
```
- Carpenter will create a new file named simple.html. Click on this file to view in your browser. The result should look something like this
TODO: insert pic

