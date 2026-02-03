# Image Organizer

This Go CLI scans a directory of images and moves each file either into a `location-map` bucket (when the filename contains the keyword) or into a `map` bucket for everything else. Defaults are already wired up for your workflow, and the flags let you tweak the keyword, destination names, or run in dry mode before touching your files.

## Default configuration

- Source folder: `C:\Users\dell\Downloads\Maps-20260203T100443Z-3-001\Maps`  
- Destination root: `C:\Users\dell\Downloads\locators images`  
- Keyword: `locator` (case-insensitive)  
- Location subfolder: `location-map`  
- Map subfolder: `map`

## Flags

- `-source` overrides the folder that is scanned for images.  
- `-dest` sets the base directory that will receive the grouped subfolders.  
- `-keyword` controls what triggers a file to be treated as a location map.  
- `-location-folder` / `-map-folder` rename the respective output directories.  
- `-dry-run` prints the planned moves without moving anything; useful for validation.  
- `-verbose` logs every move (use `-verbose=false` to quiet it).

## Example

```bash
go run . -dry-run
```

With no overrides, the program will look for images under the source path above, evaluate each filename for `locator`, and report whether it would move into `locators images/location-map` or `locators images/map`. Remove `-dry-run` to perform the copy.

## Notes

- Supported extensions: `.jpg`, `.jpeg`, `.png`, `.gif`, `.bmp`, `.tiff`, `.webp`.  
- If a name collision occurs, the program appends `_1`, `_2`, etc., to keep both files.  
- The destination root is skipped during the walk so files already moved don't get processed again.
