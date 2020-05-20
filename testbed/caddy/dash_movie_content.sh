urls=(http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/10_sec/x264/bbb/DASH_Files/live/)

folders=(html)

# get number of movies
mLen=${#urls[@]}


for (( j=0; j<${mLen}; j++ ))
do
	mkdir -p ${folders[j]}
        cd ${folders[j]}
        wget -r -np -nH --cut-dirs=3 -R index.html ${urls[j]}
        cd ..
done

