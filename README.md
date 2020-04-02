# Go CDA

Go CDA is a Clinical Document Architecture parsing library for the Go programming language. CDA's are used extensively in health care. This library was developed to parse Pathology reports saved in the CDA format.

This library is still a work in progress and the documentation will be extended with more examples in due course.

# Install

``` bash
go get github.com/PaulBradley/go-cda
```

# Example Usage

``` go
// parse the CDA document
clinicalDocument, err = cda.Parse(string(data))
if err != nil {
    log.Fatal(err.Error())
}
```

# Example Output

Once the CDA document has been parsed, it's easy to output the discrete data items in any format. Typically this would be a HTML report.

```
+------------------+---------------------+
|    DATA ITEM     |        VALUE        |
+------------------+---------------------+
| Accession #      | C12439170           |
+------------------+---------------------+
| Status Code      | completed           |
+------------------+---------------------+
| Document Title   | Laboratory Report   |
+------------------+---------------------+
| Custodian        | Holby City Hospital |
+------------------+---------------------+
| Report Date/Time | 20190814141433+0100 |
+------------------+---------------------+
| Language Code    | en-GB               |
+------------------+---------------------+
| Given Name       | Hadley              |
+------------------+---------------------+
| Surname          | Bradley             |
+------------------+---------------------+
| DOB              | 19851110            |
+------------------+---------------------+
| Gender           | Male                |
+------------------+---------------------+
```