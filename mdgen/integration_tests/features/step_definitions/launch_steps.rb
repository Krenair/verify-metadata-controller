Given("the java app exists") do
  # pending # Write code here that turns the phrase above into concrete actions
  # /integration_tests
  system('rm -fr mdgen')
  system('unzip ../mdgen.zip')
  assert File.file?("mdgen/bin/mdgen")
end

When("I run the java executable with no parameters") do
  # pending # Write code here that turns the phrase above into concrete actions
  @last_output = `mdgen/bin/mdgen 2>&1`
  # puts ("FFFFOOOOOOOOOO" + @last_output)
  # proxy ../proxy.yaml ../test/key.rsa.pem test/key.rsa.pem rsa
end

When("I run the java executable with some appropriate parameters") do
  @last_output = `java -classpath '/opt/cloudhsm/java/*:mdgen/lib/*' uk.gov.ida.mdgen.MetadataGenerator proxy ../test/proxy.yml ../test/cert.rsa.pem ../test/cert.rsa.pem --algorithm rsa 2>&1`
end

Then("I see that the application complains about missing parameters") do
  puts("HERE IS THE OUTPUT: " + @last_output)
  assert(@last_output.include? "required parameters")
end
