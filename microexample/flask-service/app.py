from flask import Flask
import requests

app = Flask(__name__)

@app.route('/countenroll')
def countenroll():

	url = "https://studentsubject-d4fe.restdb.io/rest/subject"

	headers = {
		'content-type': "application/json",
		'x-apikey': "e2ef39ac6fc3ac51270e07553ce2ba3f9b83f",
		'cache-control': "no-cache"
	}

	response = requests.request("GET", url, headers=headers)

	#print(response.text)
	resptext = str(response.text)
        a_list = []
	while resptext.find("subject\":\"") > 0:
	   reslen = len(resptext)
	   start = resptext.find("subject\":\"") + len("subject\":\"")
	   end = resptext.find("\"}")
	   substring = resptext[start:end]
           a_list.append(substring)
	   resptext = resptext[end+1:reslen]
                    
	my_dict = {i:a_list.count(i) for i in a_list}

	html = """\
	<h3> Subjectwise Enrollments</h3>
	<table border="1" cellspacing="10">
	<tr><th>Subject</th><th>Number of Enrollments</th></tr>"""
	

	for x in my_dict:	
	  html = html + "<tr>"	
	  html = html + "<td>" + str(x) + "</td>"
	  html = html + "<td>" + str(my_dict[x]) + "</td>"
	  html = html + "</tr>"

	html = html + "</table>"
	return html

if __name__ == "__main__":
    app.run(host="0.0.0.0", debug=True)
