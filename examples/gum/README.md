# Gum

### File

<img width="600" src="file.gif" />

```
Output examples/gum/file.gif

Type "gum file ./src"
Sleep 0.5s
Enter
Sleep 0.5s
Down@250ms 4
Up@250ms 1
Sleep 0.5s
Enter
Sleep 0.5s
Down@250ms 4
Sleep 0.5s
Up@250ms 2
Sleep 1s
```

### Pager

<img width="600" src="pager.gif" />

```
Output examples/gum/pager.gif

Set Padding 32
Set FontSize 16
Set Height 600

Type "gum pager < ~/src/gum/README.md --border normal"
Sleep 0.5s
Enter
Sleep 1s
Down@15ms 40
Sleep 0.5s
Up@15ms 30
Sleep 0.5s
Down@15ms 20
Sleep 2s
```

### Table

<img width="600" src="pager.gif" />

```
Output examples/gum/table.gif

Type "gum table < superhero.csv -w 2,12,5,6,6,8,4,20 --height 10"
Enter
Sleep 0.5s
Down 10
Sleep 0.5s
Down 10
Sleep 0.5s
Up 10
Sleep 0.5s
Enter
Sleep 2s
```
