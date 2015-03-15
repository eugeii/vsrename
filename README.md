# vsrename

vsrename is a utility to quicky rename a set of video files with a matching subtitle file.

Quick usage of vsrename (for full options, see below):

    vsrename
        [--subext=<subtitle extension>]
        [--vidext=<video extension>]
        [--subregex=<regex pattern to identify episode in subtitle files>]
        [--vidregex=<regex pattern to identify episode in video files>]

### Examples:

Show renames without actually renaming:

    vsrename --subext="srt" --vidext="mp4" 
             --subregex=".*1x([0-9]+).*" 
             --vidregex=".*S01E([0-9]+).*"

Show renames and actually rename:

    vsrename -w --subext="srt" --vidext="mp4" 
                --subregex=".*1x([0-9]+).*" 
                --vidregex=".*S01E([0-9]+).*"

Command line options:

    -l=".": The path to the location of the video and subtitle files
    -location=".": The path to the location of the video and subtitle files
    -subext="srt": The extension of the subtitle files (without leading '.', e.g. 'srt')
    -subregex="": The regex to identify episode of each subtitle file (as a regex group)
    -vidext="mp4": The extension of the video files (without leading '.', e.g. 'mp4')
    -vidregex="": The regex to identify episode of each video file (as a regex group)
    -w=false: Actually perform the rename (shorthand)
    -write=false: Actually perform the rename

Example output:

    Found total 13 video files (*.mp4) and 13 subtitle files (*.srt).
      [*] 'FFM.S05E01.HDTV.x264-KILLERS.mp4' -> 'Fan Fiction Movie - 5x01 - Epi1.mp4'
      [*] 'FFM.S05E02.HDTV.x264-KILLERS.mp4' -> 'Fan Fiction Movie - 5x02 - Epi2.mp4'
      [*] 'FFM.S05E03.HDTV.x264-ASAP.mp4' -> 'Fan Fiction Movie - 5x03 - Epi3.mp4'
      [*] 'FFM.S05E04.HDTV.x264-KILLERS.mp4' -> 'Fan Fiction Movie - 5x04 - Epi4.mp4'
      [*] 'FFM.S05E05.HDTV.x264-KILLERS.mp4' -> 'Fan Fiction Movie - 5x05 - Epi5.mp4'
      [*] 'FFM.S05E06.HDTV.x264-KILLERS.mp4' -> 'Fan Fiction Movie - 5x06 - Epi6.mp4'
      [*] 'FFM.S05E07.HDTV.x264-KILLERS.mp4' -> 'Fan Fiction Movie - 5x07 - Epi7.mp4'
      [*] 'FFM.S05E08.HDTV.x264-KILLERS.mp4' -> 'Fan Fiction Movie - 5x08 - Epi8.mp4'
      [*] 'FFM.S05E09.HDTV.x264-KILLERS.mp4' -> 'Fan Fiction Movie - 5x09 - Epi9.mp4'
      [*] 'FFM.S05E10.HDTV.x264-KILLERS.mp4' -> 'Fan Fiction Movie - 5x10 - Epi10.mp4'
      [*] 'FFM.S05E11.HDTV.x264-ASAP.mp4' -> 'Fan Fiction Movie - 5x11 - Epi11.mp4'
      [*] 'FFM.S05E12.HDTV.x264-KILLERS.mp4' -> 'Fan Fiction Movie - 5x12 - Epi12.mp4'
      [*] 'FFM.S05E13.HDTV.x264-KILLERS.mp4' -> 'Fan Fiction Movie - 5x13 - Epi13.mp4'

    13 files renamed.