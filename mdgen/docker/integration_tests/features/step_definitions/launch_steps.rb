AfterConfiguration() do
  system('rm -fr mdgen')
  system('unzip ../mdgen.zip')
end

Given("the java app exists") do
  expect(File.file?("mdgen/bin/mdgen")).to be true
end

Given("the login credentials exist") do
  expect(ENV.has_key? "HSM_USER").to be true
  expect(ENV.has_key? "HSM_PASSWORD").to be true
end

When("I run the java executable with no parameters") do
  @last_output = `mdgen/bin/mdgen 2>&1`
end

When("I run the java executable with some appropriate parameters") do
  @last_output = run_app("proxy", "rsa")
end

Then("I see that the application complains about missing parameters") do
  expect(@last_output).to include "required parameters"
end

private
def run_app(node_type, algorithm)
  `java -classpath '/opt/cloudhsm/java/*:mdgen/lib/*' uk.gov.ida.mdgen.MetadataGenerator #{node_type} ../test/#{node_type}.yml ../test/cert.#{algorithm}.pem ../test/cert.#{algorithm}.pem --algorithm #{algorithm} 2>&1`
end