inst_photo_downloader

A Instagram photo tool that supports concurrency download. This tool help you to download those photos for your backup, all the photos still own by original creator.

Install
--------------

  go get -u -x github.com/PavelPavells/inst_photo_downloader

Note, you need go to [Instagram developer page](https://instagram.com/developer/clients/manage/) to get latest token and update in environment variables

  export InstagramID="YOUR_ClientID_HERE"

Usage
---------------------

  inst_photo_downloader [options] 

All the photos will download to `USERS/Pictures/inst_photos`.

Options
---------------

- `-n` Instagram page name such as: [kingjames](https://instagram.com/kingjames/) 
- `-c` number of workers. (concurrency), default workers is "2"


Examples
---------------

Download all photos from LeBron James Instagram Photos with 10 workers.

  inst_photo_downloader -n=kingjames -c=10

TODOs
---------------

- Support video download.


Contribute
---------------

Please open up an issue on GitHub before you put a lot efforts on pull request.
The code submitting to PR must be filtered with `gofmt`

Related Project
---------------

Here also a Facebook Photo downloader written by Go. [https://github.com/PavelPavells/fb_photo_downloader](https://github.com/PavelPavells/fb_photo_downloader)
