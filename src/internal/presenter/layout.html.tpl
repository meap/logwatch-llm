<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Logwatch LLM - report</title>
    <style>
        body {
            font-family: Google Sans, Helvetica Neue, sans-serif;
            font-size: 16px;
            line-height: 1.8;
            margin: 0;
            background: #fff;
        }
        .container {
            max-width: 1028px;
            margin: 0 auto;
            padding: 32px 16px;
            box-sizing: border-box;
            background: #fff;
        }
        pre {
            background: #f5f5f5;
            border-radius: 8px;
            padding: 16px;
            overflow-x: auto;
        }
        code {
            background: #f5f5f5;
            border-radius: 4px;
            font-family: 'Fira Mono', 'Menlo', 'Monaco', 'Consolas', monospace;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin: 24px 0;
            background: #fafbfc;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 1px 3px rgba(0,0,0,0.04);
        }
        th, td {
            padding: 12px 16px;
            border-bottom: 1px solid #e0e0e0;
            text-align: left;
        }
        th {
            background: #f0f1f3;
            font-weight: 600;
        }
        tr:last-child td {
            border-bottom: none;
        }        
    </style>
    </head>
<body>
  <div class="container">
    {{.Content}}
  </div>
</body>
</html>
