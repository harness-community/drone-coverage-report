
1) Build the Game of Life project by running the following command in the terminal, this will install
the necessary jacoco reports jacoco.html and jacoco.xml in the target directory:
GOL_IMAGE=<add_test_data_image> sh ./build_game_of_life.sh

2) Test and check whether Game of Life image is working correctly, you should .exec, .class and .java files printed
docker run --rm -e DWPJA2=$DWPJA2 game-of-life-jacoco-image

3) Build the CheckJacocoTestDataStep2Dockerfile and follow similar steps as 1) for building and pushing docker image
  This step confirms the jacoco test data is available across all steps in the /harness directory

4) Now set up the plugin image and deploy as step 3 in Harness pipeline, consider this as step 3.
   the following values to be set in last step.

     tool: jacoco
     reports_path_pattern: "**/target/jacoco.exec"
     fail_on_threshold: "true"
     fail_if_no_reports: "false"
     class_directories: "**/target/classes, **/WEB-INF/classes"
     class_inclusion_pattern: "**/*.class, **/*.xml"
     class_exclusion_pattern: "**/controllers/*.class"
     source_directories: "**/src/main/java"
     source_inclusion_pattern: "**/*.java, *.groovy"
     source_exclusion_pattern: "**/controllers/*.java"
     skip_source_copy: "false"
     threshold_instruction: "0"
     threshold_branch: "0"
     threshold_complexity: "80"
     threshold_line: "0"
     threshold_method: "0"
     threshold_class: "0"

